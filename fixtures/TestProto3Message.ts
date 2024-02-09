// DO NOT EDIT! This file is generated automatically by pbts (github.com/octavore/pbts)

export abstract class TestProto3Message {
  strField: string;
  optStrField?: string;
  int32Field: number;
  int64Field: string;
  strList: string[];
  metadata: {[key: string]: string};
  static copy(from: TestProto3Message, to?: TestProto3Message): TestProto3Message {
    if (to) {
      to.strField = from.strField;
      to.optStrField = from.optStrField;
      to.int32Field = from.int32Field;
      to.int64Field = from.int64Field;
      to.strList = from.strList;
      to.metadata = from.metadata;
      return to;
    }
    return {...from};
  }
}
