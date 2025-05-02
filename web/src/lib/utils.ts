import { type ClassValue, clsx } from "clsx";
import { twMerge } from "tailwind-merge";

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}

export enum ButtonStateEnum {
  INIT,
  PENDING,
  SUCESS,
}

/**
 * Generates a random API key of specified length
 * @param length The length of the API key to generate (default: 24)
 * @returns A random string of the specified length
 */
export function generateRandomApiKey(length: number = 48): string {
  const characters =
    "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789";
  let result = "";
  const charactersLength = characters.length;
  for (let i = 0; i < length; i++) {
    result += characters.charAt(Math.floor(Math.random() * charactersLength));
  }
  return result;
}

/**
 * Copies text to clipboard and optionally triggers a notification
 * @param text The text to copy to clipboard
 * @param notifyCallback Optional callback function to show notifications
 * @returns A promise that resolves when the text is copied
 */
export async function copyToClipboard(
  text: string,
  notifyCallback?: ({
    title,
    content,
  }: {
    title: string;
    content: string;
  }) => void
): Promise<void> {
  try {
    await navigator.clipboard.writeText(text);
    if (notifyCallback) {
      notifyCallback({
        title: "Copied",
        content: "Content copied to clipboard",
      });
    }
  } catch (error) {
    console.error("Failed to copy: ", error);
    if (notifyCallback) {
      notifyCallback({
        title: "Error",
        content: "Failed to copy to clipboard",
      });
    }
  }
}
