// DO NOT EDIT! This file is generated automatically by pbts (github.com/octavore/pbts)

export interface TestProto3Message {
  strField: string;
  optStrField?: string;
  int32Field: number;
  int64Field: string;
  strList: string[];
  nested?: TestProto3NestedMessage;
  nestedList: TestProto3NestedMessage[];
  structField?: TestProto3GoogleStruct;
  enumField?: TestProto3EnumStruct;
  metadata: {[key: string]: number};
}

export abstract class TestProto3Message {
  static copy(from: TestProto3Message, to?: TestProto3Message): TestProto3Message {
    if (to) {
      to.strField = from.strField;
      to.optStrField = from.optStrField;
      to.int32Field = from.int32Field;
      to.int64Field = from.int64Field;
      to.strList = from.strList.slice();
      to.nested = from.nested ? TestProto3NestedMessage.copy(from.nested) : undefined;
      to.nestedList = from.nestedList.slice();
      to.structField = from.structField ? TestProto3GoogleStruct.copy(from.structField) : undefined;
      to.enumField = from.enumField ? TestProto3EnumStruct.copy(from.enumField) : undefined;
      to.metadata = structuredClone(from.metadata);
      return to;
    }
    return {
      ...from,
      nested: from.nested ? TestProto3NestedMessage.copy(from.nested) : undefined,
      nestedList: from.nestedList.slice(),
      structField: from.structField ? TestProto3GoogleStruct.copy(from.structField) : undefined,
      enumField: from.enumField ? TestProto3EnumStruct.copy(from.enumField) : undefined,
      metadata: structuredClone(from.metadata),
    };
  }
}

export interface TestProto3NestedMessage {
  strField: string;
}

export abstract class TestProto3NestedMessage {
  static copy(from: TestProto3NestedMessage, to?: TestProto3NestedMessage): TestProto3NestedMessage {
    if (to) {
      to.strField = from.strField;
      return to;
    }
    return {...from};
  }
}

export interface TestProto3GoogleStruct {
  structField?: {[key: string]: any};
}

export abstract class TestProto3GoogleStruct {
  static copy(from: TestProto3GoogleStruct, to?: TestProto3GoogleStruct): TestProto3GoogleStruct {
    if (to) {
      to.structField = structuredClone(from.structField);
      return to;
    }
    return {
      ...from,
      structField: structuredClone(from.structField),
    };
  }
}

export interface TestProto3EnumStruct {
  enumField: TestProto3EnumStruct_TestProto3Enum;
}

export abstract class TestProto3EnumStruct {
  static copy(from: TestProto3EnumStruct, to?: TestProto3EnumStruct): TestProto3EnumStruct {
    if (to) {
      to.enumField = from.enumField;
      return to;
    }
    return {...from};
  }
}

export type TestProto3EnumStruct_TestProto3Enum = 'bar' | 'foo' | 'unknown';
