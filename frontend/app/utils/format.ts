export function formatDate(input: string | Date, format = "YYYY-MM-DD HH:mm:ss"): string {
  const date = input instanceof Date ? input : new Date(input);

  if (isNaN(date.getTime())) return "";

  const map: Record<string, string> = {
    YYYY: String(date.getFullYear()),
    MM: String(date.getMonth() + 1).padStart(2, "0"),
    DD: String(date.getDate()).padStart(2, "0"),
    HH: String(date.getHours()).padStart(2, "0"),
    mm: String(date.getMinutes()).padStart(2, "0"),
    ss: String(date.getSeconds()).padStart(2, "0"),
    SSS: String(date.getMilliseconds()).padStart(3, "0"),
  };

  return Object.entries(map).reduce(
    (acc, [key, value]) => acc.replace(new RegExp(key, "g"), value),
    format,
  );
}
