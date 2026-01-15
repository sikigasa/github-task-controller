import { useState, useEffect } from "react";
import { githubApi, type GithubConnectionStatus } from "@/lib/api";

export function GitHubSettings() {
  const [status, setStatus] = useState<GithubConnectionStatus | null>(null);
  const [pat, setPat] = useState("");
  const [isLoading, setIsLoading] = useState(true);
  const [isSaving, setIsSaving] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    fetchStatus();
  }, []);

  const fetchStatus = async () => {
    try {
      const data = await githubApi.getConnectionStatus();
      setStatus(data);
    } catch (err) {
      console.error("Failed to fetch GitHub connection status:", err);
      setError("GitHub連携状態の取得に失敗しました");
    } finally {
      setIsLoading(false);
    }
  };

  const handleSavePAT = async () => {
    if (!pat.trim()) return;
    setIsSaving(true);
    setError(null);
    try {
      await githubApi.savePAT(pat);
      setPat("");
      await fetchStatus();
    } catch (err) {
      console.error("Failed to save PAT:", err);
      setError("PATの保存に失敗しました");
    } finally {
      setIsSaving(false);
    }
  };

  const handleDeletePAT = async () => {
    setIsSaving(true);
    setError(null);
    try {
      await githubApi.deletePAT();
      await fetchStatus();
    } catch (err) {
      console.error("Failed to delete PAT:", err);
      setError("PATの削除に失敗しました");
    } finally {
      setIsSaving(false);
    }
  };

  if (isLoading) {
    return <div className="p-4">読み込み中...</div>;
  }

  return (
    <div className="p-6 bg-white dark:bg-gray-800 rounded-lg shadow">
      <h2 className="text-xl font-semibold mb-4">GitHub連携設定</h2>
      
      {error && (
        <div className="mb-4 p-3 bg-red-100 text-red-700 rounded">{error}</div>
      )}

      <div className="space-y-4">
        <div className="flex items-center gap-2">
          <span className={`w-3 h-3 rounded-full ${status?.is_connected ? "bg-green-500" : "bg-gray-400"}`} />
          <span>{status?.is_connected ? `GitHub連携済み (${status.username})` : "GitHub未連携"}</span>
        </div>

        {status?.is_connected && (
          <div className="flex items-center gap-2">
            <span className={`w-3 h-3 rounded-full ${status?.has_pat ? "bg-green-500" : "bg-yellow-500"}`} />
            <span>{status?.has_pat ? "PAT設定済み" : "PAT未設定（OAuthトークン使用中）"}</span>
          </div>
        )}

        {status?.is_connected && (
          <div className="mt-6">
            <h3 className="font-medium mb-2">Personal Access Token (PAT)</h3>
            <p className="text-sm text-gray-600 dark:text-gray-400 mb-3">
              GitHub Projects V2 APIを使用するにはPATが必要です。
              <a 
                href="https://github.com/settings/tokens/new?scopes=project,read:user" 
                target="_blank" 
                rel="noopener noreferrer"
                className="text-blue-500 hover:underline ml-1"
              >
                PATを作成
              </a>
            </p>
            
            <div className="flex gap-2">
              <input
                type="password"
                value={pat}
                onChange={(e) => setPat(e.target.value)}
                placeholder="ghp_xxxxxxxxxxxx"
                className="flex-1 px-3 py-2 border rounded dark:bg-gray-700 dark:border-gray-600"
              />
              <button
                onClick={handleSavePAT}
                disabled={isSaving || !pat.trim()}
                className="px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600 disabled:opacity-50"
              >
                保存
              </button>
            </div>

            {status?.has_pat && (
              <button
                onClick={handleDeletePAT}
                disabled={isSaving}
                className="mt-2 text-sm text-red-500 hover:underline"
              >
                PATを削除
              </button>
            )}
          </div>
        )}

        {!status?.is_connected && (
          <p className="text-sm text-gray-600 dark:text-gray-400">
            GitHub連携を使用するには、まずGitHubでログインしてください。
          </p>
        )}
      </div>
    </div>
  );
}
