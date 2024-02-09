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

export enum TestEnumStruct_TestEnum {
  Bar = "bar",
  Foo = "foo",
  Unknown = "unknown",
}
