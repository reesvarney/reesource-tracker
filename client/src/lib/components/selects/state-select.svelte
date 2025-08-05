<script lang="ts">
  import * as Select from "../ui/select";
  export let bindValue: string;
  export let disabled: boolean = false;
  export let placeholder: string = "Select state";
  export let id: string = "state-select";
  export let required: boolean = false;
  export let options: { value: string; label: string }[] = [
    { value: "in_use", label: "In Use - Sample is actively being used" },
    { value: "available", label: "Available - Sample is available for use" },
    { value: "destroyed", label: "Destroyed - Sample no longer exists" },
    { value: "broken", label: "Broken - Sample does not function as required" },
    {
      value: "unassigned",
      label: "Unassigned - This code is not assigned to a sample",
    },
  ];
  export let filterMode: boolean = false;
  if (filterMode) {
    options.unshift({
      value: "",
      label: "Any",
    });
    options = options.filter((option) => option.value !== "unassigned");
  }
</script>

<Select.Root bind:value={bindValue} type="single" {disabled} {required}>
  <Select.Trigger {id} class="w-full">
    {bindValue
      ? options.find((o) => o.value === bindValue)?.label
      : placeholder}
  </Select.Trigger>
  <Select.Content class="w-full">
    {#each options as option}
      <Select.Item value={option.value}>{option.label}</Select.Item>
    {/each}
  </Select.Content>
</Select.Root>
