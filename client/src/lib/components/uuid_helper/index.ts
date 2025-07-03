export function UUIDBlobToString(sampleId: string): string {
  if (sampleId.length !== 8) {
    return sampleId;
  }
  const chars = sampleId.toUpperCase().slice(0, 6);
  // Convert base 64 to raw bytes
  if (chars.includes("=") || chars.includes("-")) {
    return sampleId; // Invalid format, return as is
  }
  const base64 = atob(chars);
  const bytes = new Uint8Array(base64.length);
  for (let i = 0; i < base64.length; i++) {
    bytes[i] = base64.charCodeAt(i);
  }
  // Convert raw bytes to base 36 string
  const base36 = Array.from(bytes)
    .map((b) => b.toString(36).toUpperCase().padStart(2, "0"))
    .join("");
  return `${base36.slice(0, 2)}-${base36.slice(2, 4)}-${base36.slice(4, 6)}`;
}

export function UUIDStringToBlob(sampleId: string): string {
  if (sampleId.length !== 8) {
    return sampleId;
  }
  const base36 = sampleId.toUpperCase().replace(/-/g, "");
  if (base36.length !== 6 || /[^0-9A-Z]/.test(base36)) {
    return sampleId; // Invalid format, return as is
  }
  // Convert base 36 string to raw bytes
  const bytes = new Uint8Array(3);
  for (let i = 0; i < 3; i++) {
    bytes[i] = parseInt(base36.slice(i * 2, i * 2 + 2), 36);
  }
  // Convert raw bytes to base64 string
  return btoa(String.fromCharCode(...bytes));
}
