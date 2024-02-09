// DO NOT EDIT! This file is generated automatically by pbts (github.com/octavore/pbts)

export abstract class TestOneofStruct {
  stock?: TestOneofStruct_Stock; // oneof:TestOneofStruct_instrument
  currency?: TestOneofStruct_Currency; // oneof:TestOneofStruct_instrument
  strField?: string; // oneof:TestOneofStruct_instrument
  static copy(from: TestOneofStruct, to?: TestOneofStruct): TestOneofStruct {
    to = to || {};
    if ('stock' in from) {
      to.stock = TestOneofStruct_Stock.copy(from.stock || {}, to.stock || {});
    }
    if ('currency' in from) {
      to.currency = TestOneofStruct_Currency.copy(from.currency || {}, to.currency || {});
    }
    to.strField = from.strField;
    return to;
  }
}

export abstract class TestOneofStruct_Stock {
  name?: string;
  static copy(from: TestOneofStruct_Stock, to?: TestOneofStruct_Stock): TestOneofStruct_Stock {
    to = to || {};
    to.name = from.name;
    return to;
  }
}

export abstract class TestOneofStruct_Currency {
  country?: string;
  shortCode?: string;
  static copy(from: TestOneofStruct_Currency, to?: TestOneofStruct_Currency): TestOneofStruct_Currency {
    to = to || {};
    to.country = from.country;
    to.shortCode = from.shortCode;
    return to;
  }
}

