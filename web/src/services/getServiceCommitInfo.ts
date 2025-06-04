import { callApi } from "./api";

export type ServiceCommitInfo = {
  hash: string;
  message: string;
  author: string;
  date: string;
  githubUrl: string;
};

export async function getServiceCommitInfo(
  serviceId: string
): Promise<ServiceCommitInfo | null> {
  try {
    return await callApi<ServiceCommitInfo>(
      `/service/${serviceId}/commit-info`,
      "GET"
    );
  } catch (error) {
    // Si l'API retourne 204 (No Content), cela signifie que ce n'est pas un service GitHub
    if (error instanceof Error && error.message.includes("204")) {
      return null;
    }
    throw error;
  }
}
