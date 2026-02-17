// DO NOT EDIT! This file is generated automatically by pbts (github.com/octavore/pbts)

export interface Struct {
  fields: {[key: string]: any};
}

export abstract class Struct {
  static copy(from: Struct, to?: Struct): Struct {
    if (to) {
      to.fields = structuredClone(from.fields);
      return to;
    }
    return {
      ...from,
      fields: structuredClone(from.fields),
    };
  }
}
