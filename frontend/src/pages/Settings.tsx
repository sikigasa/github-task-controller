import React, { useState, useEffect } from "react";
import { Github, FileText, Save, Settings as SettingsIcon } from "lucide-react";
import { Button } from "@/components/common/Button";
import { useProjects } from "@/contexts";
import { githubApi, type GithubConnectionStatus } from "@/lib/api";

export const Settings: React.FC = () => {
  const { projects } = useProjects();
  const [defaultProjectId, setDefaultProjectId] = useState("1");
  const [githubStatus, setGithubStatus] = useState<GithubConnectionStatus | null>(null);
  const [pat, setPat] = useState("");
  const [isSavingPat, setIsSavingPat] = useState(false);
  const [patError, setPatError] = useState<string | null>(null);

  useEffect(() => {
    fetchGithubStatus();
  }, []);

  const fetchGithubStatus = async () => {
    try {
      const status = await githubApi.getConnectionStatus();
      setGithubStatus(status);
    } catch (err) {
      console.error("Failed to fetch GitHub status:", err);
    }
  };

  const handleSavePat = async () => {
    if (!pat.trim()) return;
    setIsSavingPat(true);
    setPatError(null);
    try {
      await githubApi.savePAT(pat);
      setPat("");
      await fetchGithubStatus();
    } catch (err) {
      setPatError("PATの保存に失敗しました");
    } finally {
      setIsSavingPat(false);
    }
  };

  const handleDeletePat = async () => {
    setIsSavingPat(true);
    setPatError(null);
    try {
      await githubApi.deletePAT();
      await fetchGithubStatus();
    } catch (err) {
      setPatError("PATの削除に失敗しました");
    } finally {
      setIsSavingPat(false);
    }
  };

  return (
    <div className="max-w-4xl mx-auto space-y-8 pb-10">
      <div>
        <h2 className="text-2xl font-bold tracking-tight">
          Integrations & Settings
        </h2>
        <p className="text-muted-foreground">
          Manage your connections to external tools and application preferences.
        </p>
      </div>

      <div className="grid gap-6">
        {/* General Settings */}
        <div className="bg-card rounded-lg border border-border p-4 md:p-6 shadow-sm">
          <div className="flex items-start gap-4 mb-4">
            <div className="p-2 bg-primary/10 text-primary rounded-lg border border-primary/20">
              <SettingsIcon className="w-6 h-6" />
            </div>
            <div>
              <h3 className="text-lg font-semibold">General Settings</h3>
              <p className="text-sm text-muted-foreground">
                Customize your task management experience.
              </p>
            </div>
          </div>
          <div className="grid gap-4 max-w-lg ml-0 md:ml-14">
            <div className="grid gap-2">
              <label htmlFor="default-project" className="text-sm font-medium">
                Default Project
              </label>
              <p className="text-xs text-muted-foreground">
                Tasks created without a specific project will be added here.
              </p>
              <select
                id="default-project"
                value={defaultProjectId}
                onChange={(e) => setDefaultProjectId(e.target.value)}
                className="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2"
              >
                {projects.map((p) => (
                  <option key={p.id} value={p.id}>
                    {p.name}
                  </option>
                ))}
              </select>
            </div>
          </div>
        </div>

        {/* GitHub Integration */}
        <div className="bg-card rounded-lg border border-border p-4 md:p-6 shadow-sm">
          <div className="flex items-start gap-4 mb-4">
            <div className="p-2 bg-slate-900 rounded-lg text-white">
              <Github className="w-6 h-6" />
            </div>
            <div>
              <h3 className="text-lg font-semibold">GitHub</h3>
              <p className="text-sm text-muted-foreground">
                GitHub Projects V2と連携してタスクを同期します。
              </p>
            </div>
          </div>
          <div className="grid gap-4 max-w-lg ml-0 md:ml-14">
            {/* 連携状態 */}
            <div className="flex items-center gap-2">
              <span className={`w-3 h-3 rounded-full ${githubStatus?.is_connected ? "bg-green-500" : "bg-gray-400"}`} />
              <span className="text-sm">
                {githubStatus?.is_connected 
                  ? `GitHub連携済み (${githubStatus.username})` 
                  : "GitHub未連携 - GitHubでログインしてください"}
              </span>
            </div>

            {githubStatus?.is_connected && (
              <>
                <div className="flex items-center gap-2">
                  <span className={`w-3 h-3 rounded-full ${githubStatus?.has_pat ? "bg-green-500" : "bg-yellow-500"}`} />
                  <span className="text-sm">
                    {githubStatus?.has_pat ? "PAT設定済み" : "PAT未設定（Projects V2 APIには必要）"}
                  </span>
                </div>

                {patError && (
                  <div className="p-2 bg-red-100 text-red-700 rounded text-sm">{patError}</div>
                )}

                <div className="grid gap-2">
                  <label htmlFor="gh-token" className="text-sm font-medium">
                    Personal Access Token
                  </label>
                  <p className="text-xs text-muted-foreground">
                    <a 
                      href="https://github.com/settings/tokens/new?scopes=project,read:user" 
                      target="_blank" 
                      rel="noopener noreferrer"
                      className="text-blue-500 hover:underline"
                    >
                      PATを作成（project, read:user スコープが必要）
                    </a>
                  </p>
                  <div className="flex gap-2">
                    <input
                      id="gh-token"
                      type="password"
                      value={pat}
                      onChange={(e) => setPat(e.target.value)}
                      placeholder="ghp_..."
                      className="flex h-10 flex-1 rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2"
                    />
                    <Button 
                      onClick={handleSavePat} 
                      disabled={isSavingPat || !pat.trim()}
                    >
                      保存
                    </Button>
                  </div>
                  {githubStatus?.has_pat && (
                    <button
                      onClick={handleDeletePat}
                      disabled={isSavingPat}
                      className="text-sm text-red-500 hover:underline text-left"
                    >
                      PATを削除
                    </button>
                  )}
                </div>
              </>
            )}
          </div>
        </div>

        {/* Notion Integration */}
        <div className="bg-card rounded-lg border border-border p-4 md:p-6 shadow-sm">
          <div className="flex items-start gap-4 mb-4">
            <div className="p-2 bg-white border border-gray-200 rounded-lg text-black">
              <FileText className="w-6 h-6" />
            </div>
            <div>
              <h3 className="text-lg font-semibold">Notion</h3>
              <p className="text-sm text-muted-foreground">
                Import databases as task lists.
              </p>
            </div>
          </div>
          <div className="grid gap-4 max-w-lg ml-0 md:ml-14">
            <div className="grid gap-2">
              <label htmlFor="notion-token" className="text-sm font-medium">
                Integration Token
              </label>
              <input
                id="notion-token"
                type="password"
                placeholder="secret_..."
                className="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2"
              />
            </div>
            <div className="grid gap-2">
              <label htmlFor="notion-db" className="text-sm font-medium">
                Database ID
              </label>
              <input
                id="notion-db"
                type="text"
                placeholder="Notion Database ID"
                className="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2"
              />
            </div>
          </div>
        </div>

        {/* Save Button */}
        <div className="flex justify-end pt-4">
          <Button size="lg" icon={<Save className="w-4 h-4" />}>
            Save Changes
          </Button>
        </div>
      </div>
    </div>
  );
};
