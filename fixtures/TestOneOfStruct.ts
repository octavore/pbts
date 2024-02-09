// DO NOT EDIT! This file is generated automatically by pbts (github.com/octavore/pbts)

export interface TestOneofStruct {
  stock?: TestOneofStruct_Stock; // oneof:TestOneofStruct_instrument
  currency?: TestOneofStruct_Currency; // oneof:TestOneofStruct_instrument
  strField?: string; // oneof:TestOneofStruct_instrument
}

export abstract class TestOneofStruct {
  static copy(from: TestOneofStruct, to?: TestOneofStruct): TestOneofStruct {
    if (to) {
      to.stock = from.stock ? TestOneofStruct_Stock.copy(from.stock) : undefined;
      to.currency = from.currency ? TestOneofStruct_Currency.copy(from.currency) : undefined;
      to.strField = from.strField;
      return to;
    }
    return {
      ...from,
      stock: from.stock ? TestOneofStruct_Stock.copy(from.stock) : undefined,
      currency: from.currency ? TestOneofStruct_Currency.copy(from.currency) : undefined,
    }
  }
}

export interface TestOneofStruct_Stock {
  name?: string;
}

export abstract class TestOneofStruct_Stock {
  static copy(from: TestOneofStruct_Stock, to?: TestOneofStruct_Stock): TestOneofStruct_Stock {
    if (to) {
      to.name = from.name;
      return to;
    }
    return {...from};
  }
}

export interface TestOneofStruct_Currency {
  country?: string;
  shortCode?: string;
}

export abstract class TestOneofStruct_Currency {
  static copy(from: TestOneofStruct_Currency, to?: TestOneofStruct_Currency): TestOneofStruct_Currency {
    if (to) {
      to.country = from.country;
      to.shortCode = from.shortCode;
      return to;
    }
    return {...from};
  }
}
