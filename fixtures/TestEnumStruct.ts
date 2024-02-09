// DO NOT EDIT! This file is generated automatically by pbts (github.com/octavore/pbts)

export interface TestEnumStruct {
  enumField?: TestEnumStruct_TestEnum;
}

export abstract class TestEnumStruct {
  static copy(from: TestEnumStruct, to?: TestEnumStruct): TestEnumStruct {
    if (to) {
      to.enumField = from.enumField;
      return to;
    }
    return {...from};
  }
}

export type TestEnumStruct_TestEnum = 'bar' | 'foo' | 'unknown';
