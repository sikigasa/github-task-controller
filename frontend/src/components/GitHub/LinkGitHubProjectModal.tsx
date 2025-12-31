import { useState, useEffect } from "react";
import { githubApi, type GithubProject } from "@/lib/api";

interface Props {
  projectId: string;
  isOpen: boolean;
  onClose: () => void;
  onLinked: () => void;
}

export function LinkGitHubProjectModal({ projectId, isOpen, onClose, onLinked }: Props) {
  const [githubProjects, setGithubProjects] = useState<GithubProject[]>([]);
  const [selectedProject, setSelectedProject] = useState<GithubProject | null>(null);
  const [owner, setOwner] = useState("");
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (isOpen) {
      fetchGithubProjects();
    }
  }, [isOpen]);

  const fetchGithubProjects = async () => {
    setIsLoading(true);
    setError(null);
    try {
      const projects = await githubApi.listProjects();
      setGithubProjects(projects || []);
    } catch (err) {
      setError("GitHub Projectsの取得に失敗しました。PATが設定されているか確認してください。");
    } finally {
      setIsLoading(false);
    }
  };

  const handleLink = async () => {
    if (!selectedProject || !owner.trim()) return;
    
    setIsLoading(true);
    setError(null);
    try {
      await githubApi.linkProject(projectId, {
        github_owner: owner,
        github_project_number: selectedProject.number,
      });
      onLinked();
      onClose();
    } catch (err) {
      setError("連携に失敗しました");
    } finally {
      setIsLoading(false);
    }
  };

  if (!isOpen) return null;

  return (
    <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50">
      <div className="bg-white dark:bg-gray-800 rounded-lg p-6 w-full max-w-md">
        <h2 className="text-xl font-semibold mb-4">GitHub Project連携</h2>

        {error && (
          <div className="mb-4 p-3 bg-red-100 text-red-700 rounded text-sm">{error}</div>
        )}

        <div className="space-y-4">
          <div>
            <label className="block text-sm font-medium mb-1">GitHubユーザー名</label>
            <input
              type="text"
              value={owner}
              onChange={(e) => setOwner(e.target.value)}
              placeholder="your-username"
              className="w-full px-3 py-2 border rounded dark:bg-gray-700 dark:border-gray-600"
            />
          </div>

          <div>
            <label className="block text-sm font-medium mb-1">GitHub Project</label>
            {isLoading ? (
              <div className="text-sm text-gray-500">読み込み中...</div>
            ) : githubProjects.length === 0 ? (
              <div className="text-sm text-gray-500">
                Projectが見つかりません。PATの権限を確認してください。
              </div>
            ) : (
              <select
                value={selectedProject?.id || ""}
                onChange={(e) => {
                  const proj = githubProjects.find(p => p.id === e.target.value);
                  setSelectedProject(proj || null);
                }}
                className="w-full px-3 py-2 border rounded dark:bg-gray-700 dark:border-gray-600"
              >
                <option value="">選択してください</option>
                {githubProjects.map((proj) => (
                  <option key={proj.id} value={proj.id}>
                    #{proj.number} - {proj.title}
                  </option>
                ))}
              </select>
            )}
          </div>
        </div>

        <div className="flex justify-end gap-2 mt-6">
          <button
            onClick={onClose}
            className="px-4 py-2 text-gray-600 hover:bg-gray-100 rounded dark:text-gray-300 dark:hover:bg-gray-700"
          >
            キャンセル
          </button>
          <button
            onClick={handleLink}
            disabled={isLoading || !selectedProject || !owner.trim()}
            className="px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600 disabled:opacity-50"
          >
            連携
          </button>
        </div>
      </div>
    </div>
  );
}
