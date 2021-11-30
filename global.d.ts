declare namespace global {
  class Go {
    constructor();

    importObject: WebAssembly.Imports;

    run(instance: WebAssembly.Instance): Promise<unknown>;
  }
}
