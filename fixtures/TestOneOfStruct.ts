// DO NOT EDIT! This file is generated automatically by pbts (github.com/octavore/pbts)

export abstract class TestOneofStruct {
  // skipped field: instrument

  // oneof types:
  currency?: TestOneofStruct_Currency;
  stock?: TestOneofStruct_Stock;
  static copy(from: TestOneofStruct, to?: TestOneofStruct): TestOneofStruct {
    to = to || {};
    if ('currency' in from) {
      to.currency = TestOneofStruct_Currency.copy(from.currency || {}, to.currency || {});
    }
    if ('stock' in from) {
      to.stock = TestOneofStruct_Stock.copy(from.stock || {}, to.stock || {});
    }
    return to;
  }
}


// oneof types
export enum TestOneofStruct_InstrumentOneOf {
  Currency = 'currency',
  Stock = 'stock',
}
