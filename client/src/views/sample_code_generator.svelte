<script lang="ts">
  import { onMount } from "svelte";
  import QrCode from "$lib/components/qrcode/qrcode.svelte";
  import { Button } from "$lib/components/ui/button";
  import { Input } from "$lib/components/ui/input";
  import * as Card from "$lib/components/ui/card";
  import * as Alert from "$lib/components/ui/alert";
  import { CircleAlertIcon } from "lucide-svelte";
  import LocationSelect from "$lib/components/selects/location-select.svelte";
  import ProductSelect from "$lib/components/selects/product-select.svelte";
  let numSamples = $state(10);
  let generating = $state(false);
  let error = $state("");
  let printouts: { timestamp: number; ids: string[] }[] = $state([]);
  let currentSheet: string[] = $state([]);
  let showHistory = $state(false);
  let showPrintOverlay = $state(false);

  // Sizing parameters (modifiable by user)
  let qrPadding = $state(2); // mm, now used as grid gap
  let qrSize = $state(11); // mm (QR code SVG size)

  // Derived sizes
  // Sheet size in mm (A4 printable area, e.g. 200mm x 120mm)
  let sheetWidth = 198; // mm (A4 minus margins)
  let sheetHeight = 111.5; // mm

  // The grid slot size is determined by qrSize only; height for text is handled in CSS
  let gridSlotSize = $derived(qrSize); // mm
  let gridSlotPx = $derived(`${gridSlotSize}mm`);
  let qrSvgPx = $derived(`${qrSize}mm`);
  let qrPadPx = $derived(`${qrPadding}mm`); // for gap only

  // Dynamically calculate columns and rows based on user sizing (runes mode)
  let gridCols = $derived(
    Math.max(1, Math.floor((sheetWidth + qrPadding) / (qrSize + qrPadding)))
  );
  let gridRows = $derived(
    Math.max(1, Math.floor((sheetHeight + qrPadding) / (qrSize + qrPadding)))
  );
  let perSheet = $derived(gridCols * gridRows);

  onMount(() => {
    const cached = localStorage.getItem("sample_printouts");
    if (cached) {
      printouts = JSON.parse(cached);
    }
  });

  function savePrintouts() {
    // Only keep the 10 most recent printouts
    if (printouts.length > 10) printouts = printouts.slice(0, 10);
    localStorage.setItem("sample_printouts", JSON.stringify(printouts));
  }

  function GetSampleUrl(id: string) {
    return `${location.protocol}//${location.host}/app?sample_id=${id}`;
  }

  async function generateSamples(e: SubmitEvent) {
    e.preventDefault(); // Prevent default form submission
    error = "";
    generating = true;
    const res = await fetch(`/api/generate_samples?num_samples=${numSamples}`);
    if (!res.ok) {
      throw new Error("Failed to generate samples");
    }
    const data = await res.json();
    if (!data.sample_ids) throw new Error("No sample IDs returned");
    currentSheet = data.sample_ids;
    printouts = [
      { timestamp: Date.now(), ids: [...currentSheet] },
      ...printouts,
    ];
    // Only keep the 10 most recent printouts
    if (printouts.length > 10) printouts = printouts.slice(0, 10);
    savePrintouts();
    generating = false;
  }

  function SelectSheed(ids: string[]) {
    currentSheet = ids;
    showPrintOverlay = false; // Hide print overlay if showing
  }

  function deletePrintout(idx: number) {
    printouts = printouts.filter((_, i) => i !== idx);
    savePrintouts();
  }

  function printSheet(ids: string[]) {
    currentSheet = ids;
    showPrintOverlay = true;
    // Wait for overlay to render, then print
    setTimeout(() => {
      window.print();
      // Hide overlay after print (with a delay for print dialog to close)
      setTimeout(() => {
        showPrintOverlay = false;
      }, 500);
    }, 50);
  }

  function formatDate(ts: number) {
    return new Date(ts).toLocaleString();
  }

  // Helper to chunk an array into arrays of size n
  function chunkArray<T>(arr: T[], n: number): T[][] {
    const result: T[][] = [];
    for (let i = 0; i < arr.length; i += n) {
      result.push(arr.slice(i, i + n));
    }
    return result;
  }
</script>

<div class="flex flex-col gap-6">
  <Card.Root>
    <Card.Header>
      <Card.Title>Previous Code Sheets</Card.Title>
      <Alert.Root variant="destructive">
        <CircleAlertIcon class="size-4" />
        <Alert.Title>Warning</Alert.Title>
        <Alert.Description
          >Be careful to not reprint any existing codes, as this would result in
          2 samples sharing the same entry.</Alert.Description
        >
      </Alert.Root>
    </Card.Header>
    <Card.Content>
      <div class="flex flex-row max-w-full overflow-x-auto gap-8">
        {#each printouts as p, i}
          <Card.Root class="w-[32em]">
            <Card.Header>
              <Card.Title>
                {formatDate(p.timestamp)} -
                <span>{p.ids.length} samples</span>
              </Card.Title>
            </Card.Header>
            <Card.Content>
              <details>
                <summary>Show IDs</summary>
                <div
                  style="font-family:monospace; font-size:0.95em; word-break: break-all; max-height: 10em; overflow-y: auto;"
                >
                  {p.ids.join(", ")}
                </div>
              </details>
              <div class="flex flex-row justify-end gap-2">
                <Button onclick={() => SelectSheed(p.ids)} size="sm"
                  >Select</Button
                >
                <Button
                  variant="destructive"
                  onclick={() => deletePrintout(i)}
                  size="sm"
                  aria-label="Delete history item">Delete</Button
                >
              </div>
            </Card.Content>
          </Card.Root>
        {/each}
      </div>
    </Card.Content>
  </Card.Root>

  <Card.Root>
    <Card.Header>
      <Card.Title>Sample Code Generator</Card.Title>
    </Card.Header>
    <Card.Content>
      <form onsubmit={generateSamples}>
        <div class="flex flex-col gap-4">
          <label>
            Number of samples to generate:
            <Input
              type="number"
              min="1"
              max="1000"
              bind:value={numSamples}
              disabled={generating}
            />
          </label>
          <div>
            QR Padding (mm):
            <Input type="number" min="0" max="10" bind:value={qrPadding} />
          </div>
          <div>
            QR Size (mm):
            <Input type="number" min="5" max="50" bind:value={qrSize} />
          </div>
          <div>
            <Button type="submit" disabled={generating}>
              {generating ? "Generating..." : "Generate"}
            </Button>
          </div>
        </div>
      </form>
    </Card.Content>
  </Card.Root>
  {#if currentSheet.length}
    <Card.Root>
      <Card.Header>
        <Card.Title>Preview Output</Card.Title>
      </Card.Header>
      <Card.Content class=" max-h-[20vh] overflow-hidden">
        <div
          class="qr-sheet print-sheet-only"
          style="--grid-slot-size: {gridSlotPx}; --qr-svg-size: {qrSvgPx}; --qr-gap: {qrPadPx}; --grid-cols: {gridCols}; --grid-rows: {gridRows};"
        >
          {#each currentSheet as id, i}
            <div class="qr-item">
              {#if currentSheet[i]}
                <QrCode value={GetSampleUrl(currentSheet[i])} />
                <div class="qr-id">{currentSheet[i]}</div>
              {/if}
            </div>
          {/each}
        </div>
      </Card.Content>
      <Card.Footer>
        <Button
          onclick={() => printSheet(currentSheet)}
          style="margin-bottom:1em;">Print This Sheet</Button
        >
      </Card.Footer>
    </Card.Root>
  {/if}
</div>

{#if showPrintOverlay}
  <div class="print-overlay" aria-hidden={!showPrintOverlay}>
    {#each chunkArray(currentSheet, perSheet) as pageCodes, pageIdx}
      <div
        class="qr-sheet print-sheet-only"
        style="--grid-slot-size: {gridSlotPx}; --qr-svg-size: {qrSvgPx}; --qr-gap: {qrPadPx}; --grid-cols: {gridCols}; --grid-rows: {gridRows};{pageIdx >
        0
          ? 'page-break-before: always;'
          : ''}"
      >
        {#each Array(perSheet) as _, i}
          <div class="qr-item">
            {#if pageCodes[i]}
              <QrCode value={GetSampleUrl(pageCodes[i])} />
              <div class="qr-id">{pageCodes[i]}</div>
            {/if}
          </div>
        {/each}
      </div>
    {/each}
  </div>
{/if}

<style>
  /* Unified grid for both preview and print */
  .qr-sheet {
    display: grid;
    grid-template-columns: repeat(
      var(--grid-cols, 13),
      var(--grid-slot-size, 15mm)
    );
    grid-template-rows: repeat(
      var(--grid-rows, 7),
      calc(var(--grid-slot-size, 15mm) + 2em)
    );
    gap: var(--qr-gap, 2mm);
    margin: 2rem auto 0 auto;
    width: calc(
      var(--grid-slot-size, 15mm) * var(--grid-cols, 13) + var(--qr-gap, 2mm) *
        (var(--grid-cols, 13) - 1)
    );
    height: calc(
      (var(--grid-slot-size, 15mm) + 2em) * var(--grid-rows, 7) +
        var(--qr-gap, 2mm) * (var(--grid-rows, 7) - 1)
    );
    max-width: 100vw;
    box-sizing: border-box;
    justify-content: flex-start;
    align-items: flex-start;
    background: white;
    place-items: start;
    * {
      font-size: 2.2mm;
    }
  }
  .qr-item {
    width: 100%;
    height: 100%;
    margin: 0;
    padding: 0;
    page-break-inside: avoid;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    box-sizing: border-box;
    user-select: none;
    pointer-events: none;
    position: relative;
    overflow: hidden;
    :global(:not(.qr-item)) {
      flex: 1 1 0;
      max-width: 100%;
      max-height: 100%;
      aspect-ratio: 1;
      flex-shrink: 1;
      display: block;
    }
    .qr-id {
      font-family: monospace;
      width: 100%;
      overflow: hidden;
      text-align: center;
      box-sizing: border-box;
      padding: 0;
      color: black;
      background: none;
      break-inside: avoid;
      flex-shrink: 0;
    }
  }

  /* svg sizing now handled above in .qr-item :global(svg) */
  @media print {
    .print-overlay {
      display: flex !important;
      position: fixed;
      z-index: 99999;
      top: 0;
      left: 0;
      width: 100vw;
      min-height: 100vh;
      background: white;
      align-items: flex-start;
      justify-content: flex-start;
      flex-wrap: wrap;
      padding: 0;
      margin: 0;
      display: none;
      user-select: none;
      pointer-events: auto;
      box-sizing: border-box;
    }
    .print-sheet-only.qr-sheet {
      margin: 6mm auto 6mm auto;
      background: white;
    }

    * :not(.print-overlay *) {
      display: none !important;
      pointer-events: none !important;
    }
  }
</style>
