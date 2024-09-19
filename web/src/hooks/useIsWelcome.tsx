import { useLocation } from "react-router-dom";

export function useIsWelcome() {
  const { search } = useLocation();
  const queryParams = new URLSearchParams(search);
  const displayWelcomeModal = queryParams.get("welcome");
  return displayWelcomeModal;
}
