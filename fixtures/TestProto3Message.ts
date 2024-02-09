// DO NOT EDIT! This file is generated automatically by pbts (github.com/octavore/pbts)

export interface TestProto3Message {
  strField: string;
  optStrField?: string;
  int32Field: number;
  int64Field: string;
  strList: string[];
  nested?: TestProto3NestedMessage;
  nestedList: TestProto3NestedMessage[];
  metadata: {[key: string]: string};
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
      to.metadata = from.metadata;
      return to;
    }
    return {
      ...from,
      nested: from.nested ? TestProto3NestedMessage.copy(from.nested) : undefined,
      nestedList: from.nestedList.slice(),
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
