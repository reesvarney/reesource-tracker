<script lang="ts">
  import { onMount } from "svelte";
  import QRCode from "qrcode";
  let { value = "", size = 128 } = $props();
  let dataUrl = $state("");

  onMount(async () => {
    updateQRCode();
  });

  function updateQRCode() {
    if (value) {
      // @ts-ignore
      QRCode.toDataURL(value, {
        width: size,
        margin: 0,
        errorCorrectionLevel: "L",
      }).then((url: string) => (dataUrl = url));
    }
  }

  $effect(() => {
    updateQRCode();
  });
</script>

{#if dataUrl}
  <img alt="QR code" src={dataUrl} width={size} height={size} />
{/if}
