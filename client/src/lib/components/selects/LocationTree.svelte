<script lang="ts">
  import * as Select from "../ui/select";
  import type { SampleLocation } from "$lib/components/location";
  import LocationTree from "./LocationTree.svelte";
  export let locations: SampleLocation[] = [];
  export let filterOutIds: string[] = [];
  export let level: number = 0;
  export let filterMode: boolean = false;
</script>

{#each locations as location (location.id)}
  {#if !filterOutIds.includes(location.id) && location.name !== ""}
    {#if location.ChildLocations && location.ChildLocations.length > 0}
      <Select.Label style="padding-left: max({level * 1.5}em, 0.5em)"
        >{location.CombinedName}</Select.Label
      >
      {#if !filterOutIds.includes(`any-${location.id}`) && location.ChildLocations.length > 1 && filterMode}
        <Select.Item
          value={"any-" + location.id}
          style="padding-left: {(level + 1) * 1.5}em"
        >
          Anywhere in {location.CombinedName}
        </Select.Item>
      {/if}
      <LocationTree
        locations={location.ChildLocations}
        {filterOutIds}
        level={level + 1}
        {filterMode}
      />
    {:else}
      <Select.Item value={location.id} style="padding-left: {level * 1.5}em">
        {location.CombinedName}
      </Select.Item>
    {/if}
  {/if}
{/each}
