package pbts

import (
	"fmt"
	"io"
	"reflect"
	"strings"

	"sort"

	"github.com/golang/protobuf/proto"
)

const filePreamble = "// DO NOT EDIT! This file is generated automatically by pbts (github.com/octavore/pbts)\n"

func NewGenerator(w io.Writer) *Generator {
	return &Generator{
		out:    w,
		enums:  map[string]string{},
		oneofs: map[string][]string{},
	}
}

type Generator struct {
	models []reflect.Type
	enums  map[string]string
	oneofs map[string][]string
	out    io.Writer
}

func (g *Generator) RegisterMany(l ...interface{}) {
	for _, i := range l {
		g.Register(i)
	}
}
func (g *Generator) Register(i interface{}) {
	v := reflect.ValueOf(i).Type()
	if v.Kind() != reflect.Struct {
		panic("can only register struct types")
	}
	g.models = append(g.models, v)
}

func (g *Generator) Write() {
	g.p(0, filePreamble)
	for _, i := range g.models {
		g.convert(i)
	}

	// write enums
	sortedEnums := []string{}
	for t := range g.enums {
		sortedEnums = append(sortedEnums, t)
	}
	sort.Strings(sortedEnums)
	for _, t := range sortedEnums {
		e := g.enums[t]
		g.convertEnum(t, e)
	}

	// write oneofs
	g.writeOneofs()
}

func (g *Generator) p(indent int, s string) {
	spaces := strings.Repeat(" ", indent)
	fmt.Fprint(g.out, spaces, s, "\n")
}

func (g *Generator) convertEnum(typeName, enumName string) {
	enumMap := proto.EnumValueMap(enumName)
	enums := []string{}
	for enum := range enumMap {
		enums = append(enums, fmt.Sprintf("'%s'", enum))
	}
	if len(enums) > 0 {
		sort.Strings(enums)
		g.p(0, fmt.Sprintf("export type %s = %s;", typeName, strings.Join(enums, " | ")))
	}
}

func (g *Generator) writeOneofs() {
	if len(g.oneofs) > 0 {
		g.p(0, "")
		g.p(0, "// oneof types")
	}
	sortedOneofs := []string{}
	for t := range g.oneofs {
		sortedOneofs = append(sortedOneofs, t)
	}
	sort.Strings(sortedOneofs)
	for _, oneofName := range sortedOneofs {
		values := g.oneofs[oneofName]
		if len(values) > 0 {
			g.p(0, fmt.Sprintf("export enum %s {", oneofName))
			for i := 0; i < len(values); i++ {
				g.p(0, fmt.Sprintf("  %s = '%s',", strings.Title(values[i]), values[i]))
			}
			g.p(0, fmt.Sprintf("}"))
		}
	}
}

const typeFmt = "%s?: %s;"

// subconvertFields take a go Type and converts all the fields
// for that Type as Typescript fields. i.e. `myField?: aTypeScriptType`.
// Returns the list of field names.
func (g *Generator) subconvertFields(v reflect.Type) []annotatedField {
	fields := []annotatedField{}
	for j := 0; j < v.NumField(); j++ {
		f := v.Field(j)

		// name
		name := tsFieldname(f)
		if name == "-" {
			continue
		}

		// type
		typ, builtin := g.goTypeToTSType(f.Type, &f.Tag)
		if typ != "" {
			field := annotatedField{name: name}
			if !builtin {
				field.tsType = typ
			}
			fields = append(fields, field)
			g.p(2, fmt.Sprintf(typeFmt, name, typ))
		} else {
			g.p(2, "// skipped field: "+name)
		}
	}
	return fields
}

func (g *Generator) convert(v reflect.Type) {
	g.p(0, "export abstract class "+v.Name()+" {")
	fields := g.subconvertFields(v)

	// handle oneof fields
	sp := proto.GetProperties(v)
	if len(sp.OneofTypes) > 0 {
		g.p(2, "")
		g.p(2, "// oneof types:")

		// keys are sorted to ensure deterministic output
		keys := []string{}
		for key := range sp.OneofTypes {
			keys = append(keys, key)
		}
		sort.Strings(keys)

		for _, key := range keys {
			// store fields for typing later
			prop := sp.OneofTypes[key]
			oneOfField := v.Field(prop.Field)
			oneofName := v.Name() + "_" + oneOfField.Name + "OneOf"
			if g.oneofs[oneofName] == nil {
				g.oneofs[oneofName] = []string{}
			}
			if prop.Prop.JSONName != "" {
				key = prop.Prop.JSONName
			}
			g.oneofs[oneofName] = append(g.oneofs[oneofName], key)

			// merge oneof fields into parent
			f2 := g.subconvertFields(prop.Type.Elem())
			fields = append(fields, f2...)
		}
	}
	g.generateCopyFunction(v.Name(), fields)
	g.p(0, "}\n")
}

// helper function to extract subtags, e.g. `protobuf:"json=foo"`
func lookupSubTag(tag reflect.StructTag, tagName, subTag string) (string, bool) {
	t, ok := tag.Lookup(tagName)
	if !ok {
		return "", false
	}
	tParts := strings.Split(t, ",")
	prefix := subTag + "="
	for _, part := range tParts {
		if strings.HasPrefix(part, prefix) {
			return strings.TrimPrefix(part, prefix), true
		}
	}
	return "", false
}

// extract the field name from the field. prefers protobuf
// declared json name if it exists.
func tsFieldname(f reflect.StructField) string {
	proto, ok := lookupSubTag(f.Tag, "protobuf", "json")
	if ok {
		return proto
	}
	json, ok := f.Tag.Lookup("json")
	if ok {
		return strings.Split(json, ",")[0]
	}
	return strings.ToLower(f.Name)
}

// converts native go types to native ts types
var typeMap = map[string]string{
	"int":    "number", // todo: actually check number of bits in int
	"int32":  "number",
	"uint32": "number",
	"int64":  "string",
	"bool":   "boolean",
}

type protoEnum interface {
	EnumDescriptor() ([]byte, []int)
}

var protoEnumType = reflect.TypeOf((*protoEnum)(nil)).Elem()

// convert a go type to a TS type, and whether it was a TS builtin type or not.
// note: protobuf "oneof" is not supported
func (g *Generator) goTypeToTSType(t reflect.Type, tag *reflect.StructTag) (string, bool) {
	if tag != nil {
		// keep track of enums for later generation
		// AssignableTo is not strictly speaking necessary, rather it is a
		// helper to avoid unnecessary tag checks.
		if t.Name() != "" && t.AssignableTo(protoEnumType) {
			enum, ok := lookupSubTag(*tag, "protobuf", "enum")
			if ok {
				g.enums[t.Name()] = enum
			}
		}

		// do not generate oneof types
		if _, ok := tag.Lookup("protobuf_oneof"); ok {
			return "", false
		}
	}

	switch t.Kind() {
	case reflect.Ptr:
		return g.goTypeToTSType(t.Elem(), tag)
	case reflect.Slice:
		typ, _ := g.goTypeToTSType(t.Elem(), tag)
		if typ == "uint8" { // byte array
			return "string", true
		}
		typ += "[]"
		return typ, true
	case reflect.Struct:
		return t.Name(), false
	case reflect.Interface:
		return "any", true
	case reflect.Map:
		k, _ := g.goTypeToTSType(t.Key(), tag)
		e, _ := g.goTypeToTSType(t.Elem(), tag)
		return fmt.Sprintf("{ [key: %s]: %s; }", k, e), true
	default:
		typ := t.Name()
		if alt, ok := typeMap[typ]; ok {
			return alt, true
		}
		return typ, true
	}
}

type annotatedField struct {
	name   string
	tsType string
}

func (g *Generator) generateCopyFunction(class string, fields []annotatedField) {
	from := "from"
	if len(fields) == 0 {
		from = "_" // prevent unused variable warning
	}
	g.p(2, fmt.Sprintf("static copy(%s: %s, to?: %s): %s {", from, class, class, class))
	g.p(4, "to = to || {};")
	for _, field := range fields {
		if field.tsType == "" {
			g.p(4, fmt.Sprintf("to.%s = from.%s;", field.name, field.name))
		} else {
			g.p(4, fmt.Sprintf("if ('%s' in from) {", field.name))
			g.p(6, fmt.Sprintf("to.%s = %s.copy(from.%s || {}, to.%s || {});",
				field.name, field.tsType,
				field.name, field.name),
			)
			g.p(4, fmt.Sprintf("}"))
		}
	}
	g.p(4, "return to;")
	g.p(2, "}")
}
