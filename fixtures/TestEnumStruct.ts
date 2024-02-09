// DO NOT EDIT! This file is generated automatically by pbts (github.com/octavore/pbts)

export abstract class TestEnumStruct {
  enumField?: TestEnumStruct_TestEnum;
  static copy(from: TestEnumStruct, to?: TestEnumStruct): TestEnumStruct {
    if (to) {
      to.enumField = from.enumField;
      return to;
    }
    return {...from};
  }
}

export type TestEnumStruct_TestEnum = 'bar' | 'foo' | 'unknown';
