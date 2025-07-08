import type { AppData } from "$lib/components/app_store";
import { get } from "svelte/store";
import { Base64UUIDToString } from "$lib/components/id_helper";

export class SampleProduct {
  public id: string;
  public name: string;
  public parentProductID: string | null = null;
  private app_store: SvelteStore<AppData> | null = null;

  constructor(
    product_data: { [key: string]: any },
    app_store: SvelteStore<AppData> | null = null
  ) {
    this.id = Base64UUIDToString(product_data.ID);
    this.name = product_data.Name || "";
    this.parentProductID =
      Base64UUIDToString(product_data.ParentProductID) || null;
    this.app_store = app_store;
  }

  get ParentProduct(): SampleProduct | null {
    if (!this.parentProductID) {
      return null;
    }
    if (!this.app_store) {
      return null;
    }
    const product = get(this.app_store).products.find(
      (p) => p.id === this.parentProductID
    );
    return product ?? null;
  }

  get ChildProducts(): SampleProduct[] {
    if (!this.app_store) {
      return [];
    }
    return get(this.app_store).products.filter(
      (p) => p.parentProductID === this.id
    );
  }

  get CombinedName(): string {
    if (this.ParentProduct) {
      return `${this.ParentProduct.CombinedName} - ${this.name}`;
    }
    return this.name;
  }
}
