<script lang="ts">
  import * as Select from "../ui/select";
  export let bindValue: string;
  export let disabled: boolean = false;
  export let placeholder: string = "Select state";
  export let id: string = "state-select";
  export let required: boolean = false;
  export let options: { value: string; label: string }[] = [
    { value: "in_use", label: "In Use" },
    { value: "available", label: "Available" },
    { value: "destroyed", label: "Destroyed" },
    { value: "broken", label: "Broken" },
    { value: "unassigned", label: "Unassigned" },
  ];
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
