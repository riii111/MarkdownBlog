import type { UseFetchOptions } from "nuxt/app";

export type HttpMethod = "GET" | "POST" | "PUT" | "DELETE" | "PATCH";

interface IFetchOptions<T = unknown> {
  method: HttpMethod;
  params?: Record<string, unknown>;
  body?: T;
}

function createEndpoint(path: string) {
  // 末尾にスラッシュ追加する事でプロキシ・バックエンドなどのサーバー側のリダイレクトを回避し、パフォーマンスを僅かに向上させる
  const apiPath = path.startsWith("/") ? path : `/${path}`;
  return `/api${apiPath}${apiPath.endsWith("/") ? "" : "/"}`;
}

// エラーメッセージの取得
function getErrorMessage(status: number): string {
  switch (status) {
    case 400:
      return "リクエストに問題があります。入力内容を確認してください。";
    case 401:
      return "認証セッションの有効期限が切れました。再度ログインしてください。";
    case 403:
      return "アクセス権限がありません。";
    case 404:
      return "リソースが見つかりません。";
    case 413:
      return "アップロードされたファイルが大きすぎます。";
    case 422:
      return "入力内容に誤りがあります。";
    case 500:
      return "サーバーエラーが発生しました。";
    case 502:
      return "アクセスが集中しています。";
    default:
      return "予期せぬエラーが発生しました。";
  }
}

export function useCoreApi<ResponseType, RequestType = unknown>(
  endpoint: string,
  options: IFetchOptions<RequestType> = { method: "GET" }
) {
  const config = useRuntimeConfig();
  const fullUrl = createEndpoint(endpoint);

  const defaults: UseFetchOptions<ResponseType> = {
    baseURL: config.public.apiBase as string,
    method: options.method,
    credentials: "include",
    headers: new Headers({
      "Content-Type": "application/json",
      Accept: "application/json",
    }),
    cache: "no-store",
  };

  if (options.params && options.method === "GET") {
    defaults.params = options.params;
  }

  if (options.body) {
    if (options.body instanceof FormData) {
      (defaults.headers as Headers).delete("Content-Type");
      defaults.body = options.body;
    } else {
      defaults.body = options.body;
    }
  }

  function handleError(status: number) {
    throw createError({
      statusCode: status,
      message: getErrorMessage(status),
    });
  }

  defaults.onResponseError = (context) => {
    handleError(context.response.status);
  };

  return useFetch<ResponseType extends void ? unknown : ResponseType>(
    fullUrl,
    defaults
  );
}
