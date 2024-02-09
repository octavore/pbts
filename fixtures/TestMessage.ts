// DO NOT EDIT! This file is generated automatically by pbts (github.com/octavore/pbts)

export interface TestMessage {
  strField?: string;
  int32Field?: number;
  int64Field?: string;
  strList: string[];
  metadata: {[key: string]: string};
}

export abstract class TestMessage {
  static copy(from: TestMessage, to?: TestMessage): TestMessage {
    if (to) {
      to.strField = from.strField;
      to.int32Field = from.int32Field;
      to.int64Field = from.int64Field;
      to.strList = from.strList.slice();
      to.metadata = from.metadata;
      return to;
    }
    return {...from};
  }
}
