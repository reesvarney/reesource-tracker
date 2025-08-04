<script lang="ts">
  import QRScanner from "$lib/components/qr_scanner/qr_scanner.svelte";
  import * as InputOTP from "$lib/components/ui/input-otp";
  let { active = $bindable(false) } = $props();

  let selectedVideoInput: string = $state("");
  let rootEl: HTMLElement | null = null;
  let new_sample: string = $state("");

  $effect(() => {
    if (new_sample.length == 6) {
      window.location.assign(
        `/app?sample_id=${new_sample.slice(0, 2)}-${new_sample.slice(2, 4)}-${new_sample.slice(4, 6)}`
      );
    }
  });
  // No need to enumerate video inputs here; handled by QRScanner

  // QR scan handler
  function handleQRScan(text: string) {
    try {
      const url = new URL(text);
      // Only open if same host
      if (url.host === window.location.host) {
        window.location.href = url.href;
      } else {
        alert("Scanned QR code is not for this host.");
      }
    } catch (e) {
      alert("Scanned QR code is not a valid URL.");
    }
  }

  // No need to enumerate video inputs here; handled by QRScanner
</script>

<div class="space-y-4 flex flex-col" bind:this={rootEl}>
  <div class="flex flex-col items-center">
    <label for="qr-reader-find" class="mb-6">Scan a QR code</label>
    <QRScanner
      containerId="qr-reader-find"
      bind:selectedVideoInput
      onQrCodeScan={handleQRScan}
      autoStart={active}
    />
  </div>
  <div class="self-center mt-12 flex flex-col items-center gap-6">
    <label for="id-input">Or manually enter the sample ID</label>
    <InputOTP.Root maxlength={6} bind:value={new_sample} id="id-input">
      {#snippet children({ cells })}
        <InputOTP.Group>
          {#each cells.slice(0, 2) as cell}
            <InputOTP.Slot {cell} />
          {/each}
        </InputOTP.Group>
        <InputOTP.Separator />
        <InputOTP.Group>
          {#each cells.slice(2, 4) as cell}
            <InputOTP.Slot {cell} />
          {/each}
        </InputOTP.Group>
        <InputOTP.Separator />
        <InputOTP.Group>
          {#each cells.slice(4, 6) as cell}
            <InputOTP.Slot {cell} />
          {/each}
        </InputOTP.Group>
      {/snippet}
    </InputOTP.Root>
  </div>
</div>

<style>
  #qr-reader-find {
    width: 100%;
    max-width: 400px;
    margin: auto;
  }
</style>
