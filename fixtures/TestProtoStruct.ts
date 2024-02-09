// DO NOT EDIT! This file is generated automatically by pbts (github.com/octavore/pbts)

export abstract class Struct {
  fields?: {[key: string]: any};
  static copy(from: Struct, to?: Struct): Struct {
    to = to || {};
    to.fields = from.fields;
    return to;
  }
}

