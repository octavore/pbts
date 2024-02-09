// DO NOT EDIT! This file is generated automatically by pbts (github.com/octavore/pbts)

export abstract class TestMessage {
  strField?: string;
  int32Field?: number;
  int64Field?: string;
  strList: string[];
  metadata: {[key: string]: string};
  static copy(from: TestMessage, to?: TestMessage): TestMessage {
    if (to) {
      to.strField = from.strField;
      to.int32Field = from.int32Field;
      to.int64Field = from.int64Field;
      to.strList = from.strList;
      to.metadata = from.metadata;
      return to;
    }
    return {...from};
  }
}
