/**
 * Utility functions to manage API key storage
 */

const API_KEY_STORAGE_KEY = "api_key";

/**
 * Saves the API key to localStorage
 * @param apiKey - The API key to save
 */
export function saveApiKey(apiKey: string): void {
  try {
    localStorage.setItem(API_KEY_STORAGE_KEY, apiKey);
  } catch (error) {
    console.error("Failed to save API key to localStorage:", error);
  }
}

/**
 * Retrieves the API key from localStorage
 * @returns The stored API key or null if not found
 */
export function getStoredApiKey(): string | null {
  try {
    return localStorage.getItem(API_KEY_STORAGE_KEY);
  } catch (error) {
    console.error("Failed to retrieve API key from localStorage:", error);
    return null;
  }
}

/**
 * Removes the API key from localStorage
 */
export function removeApiKey(): void {
  try {
    localStorage.removeItem(API_KEY_STORAGE_KEY);
  } catch (error) {
    console.error("Failed to remove API key from localStorage:", error);
  }
}

/**
 * Checks if an API key exists in localStorage
 * @returns True if an API key exists, false otherwise
 */
export function hasApiKey(): boolean {
  return getStoredApiKey() !== null;
}
