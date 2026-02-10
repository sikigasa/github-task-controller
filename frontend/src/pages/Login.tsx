import { Github, Chrome, Command, Sparkles, AlertCircle } from "lucide-react";
import { useAuth } from "@/contexts";
import { useSearchParams } from "react-router-dom";

const errorMessages: Record<string, string> = {
  invalid_state: "認証セッションが無効です。もう一度お試しください。",
  no_code: "認証コードが取得できませんでした。",
  auth_failed: "認証に失敗しました。もう一度お試しください。",
  session_failed: "セッションの保存に失敗しました。",
};

export const Login: React.FC = () => {
  const { loginWithGoogle, loginWithGithub } = useAuth();
  const [searchParams] = useSearchParams();
  const error = searchParams.get("error");

  return (
    <div className="min-h-screen w-full flex items-center justify-center bg-[#020817] relative overflow-hidden text-slate-200 selection:bg-blue-500/30">
      {/* Background Effects */}
      <div className="absolute top-[-10%] left-[-10%] w-[800px] h-[800px] bg-blue-600/20 rounded-full blur-[120px] animate-pulse" />
      <div className="absolute bottom-[-10%] right-[-10%] w-[600px] h-[600px] bg-cyan-500/10 rounded-full blur-[100px] animate-pulse" />
      <div className="absolute top-[20%] right-[20%] w-[300px] h-[300px] bg-indigo-500/20 rounded-full blur-[80px]" />

      {/* Sparkle Accents */}
      <div className="absolute top-10 left-10 text-blue-400/30">
        <Sparkles className="w-8 h-8" />
      </div>
      <div className="absolute bottom-20 right-20 text-cyan-400/20">
        <Sparkles className="w-12 h-12" />
      </div>

      {/* Login Card */}
      <div className="w-full max-w-md bg-white/5 backdrop-blur-2xl border border-white/10 shadow-2xl rounded-3xl p-8 relative z-10 mx-4 overflow-hidden group">
        <div className="absolute inset-0 bg-gradient-to-br from-white/5 via-transparent to-transparent opacity-0 group-hover:opacity-100 transition-opacity duration-500 pointer-events-none" />

        <div className="flex flex-col items-center text-center space-y-8 relative z-10">
          {/* Logo */}
          <div className="relative">
            <div className="absolute inset-0 bg-blue-500 blur-xl opacity-20 rounded-full" />
            <div className="w-20 h-20 bg-gradient-to-tr from-blue-600 to-cyan-400 rounded-2xl flex items-center justify-center shadow-lg transform hover:scale-105 transition-transform duration-300 relative z-10">
              <span className="text-4xl font-black text-white tracking-tighter">
                T
              </span>
            </div>
          </div>

          <div className="space-y-2">
            <h1 className="text-3xl font-bold bg-gradient-to-r from-blue-200 via-blue-100 to-cyan-200 bg-clip-text text-transparent">
              Welcome to TaskCtl
            </h1>
            <p className="text-slate-400 text-sm">
              Experience the flow of intelligent task management.
            </p>
          </div>

          {/* OAuth Error Message */}
          {error && (
            <div className="w-full bg-red-500/10 border border-red-500/30 rounded-xl p-4 flex items-start gap-3">
              <AlertCircle className="w-5 h-5 text-red-400 flex-shrink-0 mt-0.5" />
              <p className="text-red-300 text-sm">
                {errorMessages[error] || `認証エラー: ${error}`}
              </p>
            </div>
          )}

          {/* Social Login Buttons */}
          <div className="w-full space-y-3 pt-2">
            <button
              onClick={loginWithGoogle}
              className="w-full h-12 bg-white text-slate-900 rounded-xl font-semibold flex items-center justify-center gap-3 transition-all transform hover:scale-[1.02] hover:bg-blue-50 shadow-[0_0_20px_rgba(255,255,255,0.1)] hover:shadow-[0_0_25px_rgba(59,130,246,0.3)] ring-1 ring-slate-200/50"
            >
              <Chrome className="w-5 h-5 text-blue-600" />
              <span>Continue with Google</span>
            </button>

            <button
              onClick={loginWithGithub}
              className="w-full h-12 bg-[#1e293b] hover:bg-[#334155] text-white rounded-xl font-medium flex items-center justify-center gap-3 transition-all transform hover:scale-[1.02] shadow-lg border border-slate-700"
            >
              <Github className="w-5 h-5" />
              <span>Continue with GitHub</span>
            </button>
          </div>

          {/* Divider */}
          <div className="relative w-full py-2">
            <div className="absolute inset-0 flex items-center">
              <div className="w-full border-t border-white/5" />
            </div>
            <div className="relative flex justify-center text-xs uppercase tracking-widest">
              <span className="bg-[#0b1221] px-2 text-slate-500">or</span>
            </div>
          </div>

          {/* SSO Placeholder */}
          <button className="text-sm text-slate-400 hover:text-blue-300 transition-colors flex items-center gap-2 group/sso">
            <Command className="w-4 h-4 group-hover/sso:text-blue-400 transition-colors" />
            <span>Enterprise SSO</span>
          </button>
        </div>
      </div>

      {/* Footer */}
      <div className="absolute bottom-6 flex flex-col items-center gap-2 text-[10px] text-slate-600 uppercase tracking-widest">
        <div>Task Controller V1.0</div>
        <div className="flex gap-4">
          <a href="#" className="hover:text-slate-400 transition-colors">
            Privacy
          </a>
          <a href="#" className="hover:text-slate-400 transition-colors">
            Terms
          </a>
        </div>
      </div>
    </div>
  );
};
