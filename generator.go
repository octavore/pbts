package pbts

import (
	"fmt"
	"io"
	"strings"

	"sort"

	"github.com/iancoleman/strcase"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/structpb"
)

const filePreamble = "// DO NOT EDIT! This file is generated automatically by pbts (github.com/octavore/pbts)"

func NewGenerator(w io.Writer) *Generator {
	return &Generator{out: w}
}

type Generator struct {
	models         []protoreflect.MessageDescriptor
	indirectModels []protoreflect.MessageDescriptor
	enums          []protoreflect.EnumDescriptor
	out            io.Writer

	NativeEnums bool
}

func (g *Generator) RegisterMany(l ...proto.Message) {
	for _, m := range l {
		g.Register(m)
	}
}

func (g *Generator) Register(msg proto.Message) {
	g.RegisterDescriptor(msg.ProtoReflect().Descriptor())
}

func (g *Generator) RegisterDescriptor(msgDesc protoreflect.MessageDescriptor) {
	g.models = append(g.models, msgDesc)

	for i := 0; i < msgDesc.Messages().Len(); i++ {
		m := msgDesc.Messages().Get(i)
		if m.IsMapEntry() {
			continue
		}
		g.indirectModels = append(g.indirectModels, m)
	}
}

func (g *Generator) Write() {
	g.p(0, filePreamble)
	for _, m := range g.models {
		g.p(0, "")
		g.convert(m)
	}
	for _, m := range g.indirectModels {
		g.p(0, "")
		g.convert(m)
	}

	// write enums
	g.writeEnums()
}

func (g *Generator) p(indent int, s string, args ...any) {
	spaces := strings.Repeat(" ", indent)
	fmt.Fprintf(g.out, spaces+s+"\n", args...)
}

func (g *Generator) convertEnum(e protoreflect.EnumDescriptor) {
	enumValues := e.Values()
	if g.NativeEnums {
		enums := []string{}
		for i := 0; i < enumValues.Len(); i++ {
			enum := enumValues.Get(i)
			enums = append(enums, fmt.Sprintf("%s", enum.Name()))
		}
		sort.Strings(enums)
		g.p(0, "export enum %s {", nameWithParent(e))
		for _, e := range enums {
			g.p(2, "%s = \"%s\",", strcase.ToCamel(e), e)
		}
		g.p(0, "}")
		return
	}

	enums := []string{}
	for i := 0; i < enumValues.Len(); i++ {
		enum := enumValues.Get(i)
		enums = append(enums, fmt.Sprintf("'%s'", enum.Name()))
	}

	sort.Strings(enums)
	g.p(0, "export type %s = %s;", nameWithParent(e), strings.Join(enums, " | "))
}

func (g *Generator) writeEnums() {
	sort.Slice(g.enums, func(i, j int) bool {
		return g.enums[i].Name() < g.enums[j].Name()
	})
	for _, enum := range g.enums {
		g.p(0, "")
		g.convertEnum(enum)
	}
}

// subconvertFields take a go Type and converts all the fields
// for that Type as Typescript fields. i.e. `myField?: aTypeScriptType`.
// Returns the list of field names.
func (g *Generator) subconvertFields(v protoreflect.FieldDescriptors) []annotatedField {
	fields := []annotatedField{}
	for j := 0; j < v.Len(); j++ {
		fieldDesc := v.Get(j)
		name := fieldDesc.JSONName()
		typ, builtin := g.fieldToBaseType(fieldDesc)
		if typ != "" {
			field := annotatedField{name: name, required: !fieldDesc.HasPresence(), repeated: fieldDesc.IsList()}
			if !builtin {
				field.tsType = typ
			}
			fields = append(fields, field)

			comment := ""
			if fieldDesc.ContainingOneof() != nil && !fieldDesc.HasOptionalKeyword() {
				comment = fmt.Sprintf(" // oneof:%s", nameWithParent(fieldDesc.ContainingOneof()))
			}
			if fieldDesc.HasPresence() {
				g.p(2, "%s?: %s;%s", name, typ, comment)
			} else {
				g.p(2, "%s: %s;%s", name, typ, comment)
			}
		} else {
			g.p(2, "// skipped field: "+name)
		}
	}
	return fields
}

func startsWithLower(s string) bool {
	for _, c := range s {
		if strings.ToLower(string(c)) == string(c) {
			return true
		}
		break
	}
	return false
}

func (g *Generator) convert(msg protoreflect.MessageDescriptor) {
	className := nameWithParent(msg)
	g.p(0, "export interface %s {", className)
	fields := g.subconvertFields(msg.Fields())
	g.p(0, "}\n")
	g.p(0, "export abstract class %s {", className)
	g.generateCopyFunction(className, fields)
	g.p(0, "}")
}

func (g *Generator) fieldToBaseType(f protoreflect.FieldDescriptor) (string, bool) {
	// 1. builtin kinds
	builtin := ""
	switch f.Kind() {
	case protoreflect.BoolKind:
		builtin = "boolean"
	case protoreflect.StringKind:
		builtin = "string"
	case protoreflect.Int32Kind,
		protoreflect.Sint32Kind,
		protoreflect.Uint32Kind,
		protoreflect.FloatKind:
		builtin = "number"

	case protoreflect.Int64Kind,
		protoreflect.Sint64Kind,
		protoreflect.Uint64Kind:
		builtin = "string"
	}
	if builtin != "" {
		if f.Cardinality() == protoreflect.Repeated {
			builtin += "[]"
		}
		return builtin, true
	}
	// case protoreflect.Sfixed32Kind:
	// case protoreflect.Fixed32Kind:
	// case protoreflect.Sfixed64Kind:
	// case protoreflect.Fixed64Kind:
	// case protoreflect.DoubleKind:
	// case protoreflect.BytesKind:

	// 2. enum kinds
	if f.Kind() == protoreflect.EnumKind {
		g.enums = append(g.enums, f.Enum())
		return nameWithParent(f.Enum()), true
	}

	// 3. message kinds
	if f.Kind() == protoreflect.MessageKind {
		if f.Message() == (&structpb.Value{}).ProtoReflect().Descriptor() {
			return "any", true
		}
		if f.Message() == (&anypb.Any{}).ProtoReflect().Descriptor() {
			return "any", true
		}
		if f.IsMap() {
			k, _ := g.fieldToBaseType(f.MapKey())
			v, _ := g.fieldToBaseType(f.MapValue())
			return fmt.Sprintf("{[key: %s]: %s}", k, v), true
		}
		if f.Cardinality() == protoreflect.Repeated {
			return nameWithParent(f.Message()) + "[]", false
		}
		return nameWithParent(f.Message()), false
	}

	return "", true
}

func nameWithParent(e protoreflect.Descriptor) string {
	name := e.Name()
	_, isFileParent := e.Parent().(protoreflect.FileDescriptor)
	if e.Parent() != nil && !isFileParent {
		return string(e.Parent().Name() + "_" + name)
	}
	return string(name)
}

type annotatedField struct {
	name     string
	tsType   string
	required bool
	repeated bool
}

func (a *annotatedField) generateCopiedValue() string {
	if a.repeated {
		return fmt.Sprintf("from.%s.slice()", a.name)
	}
	if a.tsType == "" || a.required {
		return fmt.Sprintf("from.%s", a.name)
	}
	return fmt.Sprintf("from.%s ? %s.copy(from.%s) : undefined", a.name, a.tsType, a.name)
}

func (g *Generator) generateCopyFunction(class string, fields []annotatedField) {
	from := "from"
	if len(fields) == 0 {
		from = "_" // prevent unused variable warning
	}
	g.p(2, "static copy(%s: %s, to?: %s): %s {", from, class, class, class)
	g.p(4, "if (to) {")
	for _, field := range fields {
		g.p(6, "to.%s = %s;", field.name, field.generateCopiedValue())
	}
	g.p(6, "return to;")
	g.p(4, "}")

	explicitCopy := []string{}
	for _, field := range fields {
		if field.tsType != "" {
			explicitCopy = append(explicitCopy,
				fmt.Sprintf("%s: %s,", field.name, field.generateCopiedValue()))
		}
	}
	if len(explicitCopy) == 0 {
		g.p(4, "return {...from};")
	} else {
		g.p(4, "return {")
		g.p(6, "...from,")
		for _, l := range explicitCopy {
			g.p(6, l)
		}
		g.p(4, "};")
	}
	g.p(2, "}")
}
