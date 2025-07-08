import { Base64UUIDToString } from "$lib/components/id_helper";
export class SampleMod {
  id: string;
  sample_id: string;
  name: string;
  time_added: Date;
  time_removed: Date | null;

  constructor(mod_data: { [key: string]: any }) {
    this.id = Base64UUIDToString(mod_data.ID);
    this.sample_id = mod_data.SampleID;
    this.name = mod_data.Name;
    this.time_added = new Date(mod_data.TimeAdded);
    if (mod_data.TimeRemoved && mod_data.TimeRemoved.Valid) {
      this.time_removed = new Date(mod_data.TimeRemoved.Time);
    } else {
      this.time_removed = null;
    }
  }
}
