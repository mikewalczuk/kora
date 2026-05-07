export const customFetch = async <T>(
  url: string,
  options: RequestInit,
): Promise<T> => {
  const baseUrl = import.meta.env.VITE_API_URL ?? "/api";
  const response = await fetch(`${baseUrl}${url}`, {
    ...options,
    credentials: "include",
  });

  if (!response.ok) throw response;

  const text = await response.text();
  return (text ? JSON.parse(text) : undefined) as T;
};
