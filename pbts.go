package pbts

import (
	"fmt"
	"os"
	"reflect"
	"sort"
	"strings"

	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
)

type messageType struct {
	Name     string
	Instance interface{}
}

type options struct {
	exclusions  []string
	nativeEnums bool
	verbose     bool
}

func GenerateAll(destPath string, optFns ...optFn) error {
	o := &options{}
	for _, optFn := range optFns {
		optFn(o)
	}

	f, err := os.OpenFile(destPath, os.O_TRUNC|os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		return fmt.Errorf("error opening file: %w", err)
	}

	g := NewGenerator(f)
	g.NativeEnums = o.nativeEnums

	messageTypes := []*messageType{}
	protoregistry.GlobalTypes.RangeMessages(func(m protoreflect.MessageType) bool {
		messageName := fmt.Sprintf("%s", m.Descriptor().FullName())
		for _, prefix := range o.exclusions {
			if strings.HasPrefix(messageName, prefix) {
				return true
			}
		}
		i := reflect.ValueOf(m.New().Interface()).Elem().Interface()
		messageTypes = append(messageTypes, &messageType{messageName, i})
		return true
	})

	sort.Slice(messageTypes, func(i, j int) bool {
		return messageTypes[i].Name < messageTypes[j].Name
	})
	for _, m := range messageTypes {
		if o.verbose {
			fmt.Println(m.Name)
		}
		g.Register(m.Instance)
	}
	g.Write()

	return nil
}

type optFn func(o *options)

func WithExclusions(exclusions ...string) optFn {
	return func(o *options) {
		o.exclusions = append(o.exclusions, exclusions...)
	}
}

func WithNativeEnums() optFn {
	return func(o *options) {
		o.nativeEnums = true
	}
}

func WithVerbose() optFn {
	return func(o *options) {
		o.verbose = true
	}
}
