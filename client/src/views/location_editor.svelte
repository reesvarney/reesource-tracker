<script lang="ts">
  import { Button } from "$lib/components/ui/button";
  import * as Select from "$lib/components/ui/select";
  import LocationSelect from "$lib/components/selects/location-select.svelte";
  import { Input } from "$lib/components/ui/input";
  import { AppStore } from "$lib/components/app_store";
  import {
    Table,
    TableHeader,
    TableHead,
    TableBody,
    TableRow,
    TableCell,
  } from "$lib/components/ui/table";

  import type { SampleLocation } from "$lib/components/location";
  import { toast } from "svelte-sonner";

  let updateTimeouts = new Map();
  // Store for parent location IDs, since it needs a string value for Select
  let parent_location_id: { [key: string]: string } = $state({});

  AppStore.subscribe((data) => {
    data.locations.forEach((location) => {
      if (location.parentLocationID) {
        parent_location_id[location.id] = location.parentLocationID;
      } else {
        parent_location_id[location.id] = "";
      }
    });
  });
  async function updateLocation(location: SampleLocation) {
    let payload = {
      name: location.name,
      description: location.description,
      parent_location_id: parent_location_id[location.id] || null,
    };
    const res = await fetch(`/api/location/${location.id}`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(payload),
    });
    if (!res.ok) {
      toast.error("Failed to update location.");
      throw new Error(await res.text());
    }
    toast.success("Location updated successfully.");
  }

  function debounceUpdate(location: SampleLocation) {
    if (updateTimeouts.has(location.id)) {
      clearTimeout(updateTimeouts.get(location.id));
    }
    const timeout = setTimeout(() => {
      updateLocation(location);
      updateTimeouts.delete(location.id);
    }, 1000);
    updateTimeouts.set(location.id, timeout);
  }

  async function addLocationRow() {
    await fetch("/api/location", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        name: "",
        description: "",
        parent_location_id: null,
      }),
    });
  }

  async function deleteLocation(location: SampleLocation) {
    await fetch(`/api/location/${location.id}`, { method: "DELETE" });
  }
</script>

<Table class="w-full">
  <TableHeader>
    <TableRow>
      <TableHead>Name</TableHead>
      <TableHead>Description</TableHead>
      <TableHead>Parent Location</TableHead>
      <TableHead>Delete</TableHead>
    </TableRow>
  </TableHeader>
  <TableBody>
    {#each $AppStore.locations as location, i (location.id || i)}
      <TableRow>
        <TableCell>
          <Input
            class="input input-bordered w-full"
            bind:value={location.name}
            oninput={(e: Event) => {
              const target = e.target as HTMLInputElement;
              location.name = target.value;
              debounceUpdate(location);
            }}
            placeholder="Location Name"
          />
        </TableCell>
        <TableCell>
          <Input
            class="input input-bordered w-full"
            bind:value={location.description}
            onchange={(e: Event) => {
              const target = e.target as HTMLInputElement;
              location.description = target.value;
              debounceUpdate(location);
            }}
            placeholder="Description"
          />
        </TableCell>
        <TableCell>
          <LocationSelect
            bind:bindValue={parent_location_id[location.id]}
            onValueChange={(v) => {
              parent_location_id[location.id] = v;
              parent_location_id = parent_location_id;
              debounceUpdate(location);
            }}
            placeholder="None"
            id={`parent-location-${location.id}`}
            disabled={false}
            required={false}
            filterOutIds={[
              location.id,
              ...location.ChildLocations.map((child) => child.id),
            ]}
            options={$AppStore.locations
              .filter(
                (l) =>
                  l.id !== location.id &&
                  !location.ChildLocations.some((child) => child.id === l.id)
              )
              .map((l) => ({ value: l.id, label: l.name }))}
          />
        </TableCell>
        <TableCell>
          <Button
            type="button"
            variant="destructive"
            onclick={() => deleteLocation(location)}
          >
            Delete
          </Button>
        </TableCell>
      </TableRow>
    {/each}
  </TableBody>
</Table>
<div class="mt-4 w-full flex flex-row justify-end">
  <Button type="button" onclick={addLocationRow}>Add Location</Button>
</div>
