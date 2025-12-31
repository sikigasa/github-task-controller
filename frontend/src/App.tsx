import { useState } from 'react';
import { Layout } from '@/components/Layout/Layout';
import { ProjectList } from '@/pages/ProjectList';
import { ProjectBoard } from '@/pages/ProjectBoard';
import { AllTasks } from '@/pages/AllTasks';
import { Settings } from '@/pages/Settings';
import { Login } from '@/pages/Login';
import { TaskProvider, ProjectProvider, AuthProvider, useAuth } from '@/contexts';

function AppContent() {
  const { isAuthenticated, isLoading, logout } = useAuth();
  const [view, setView] = useState('my-tasks');

  // ローディング中
  if (isLoading) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-background">
        <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary" />
      </div>
    );
  }

  // 未認証
  if (!isAuthenticated) {
    return <Login />;
  }

  const renderContent = () => {
    if (view === 'projects') {
      return <ProjectList onSelectProject={(id) => setView(`project:${id}`)} />;
    }
    if (view === 'my-tasks') {
      return <AllTasks />;
    }
    if (view === 'settings') {
      return <Settings />;
    }
    if (view.startsWith('project:')) {
      const projectId = view.split(':')[1];
      return <ProjectBoard projectId={projectId} />;
    }
    return <ProjectList onSelectProject={(id) => setView(`project:${id}`)} />;
  };

  return (
    <Layout activeView={view} onNavigate={setView} onLogout={logout}>
      {renderContent()}
    </Layout>
  );
}

function App() {
  return (
    <AuthProvider>
      <ProjectProvider>
        <TaskProvider>
          <AppContent />
        </TaskProvider>
      </ProjectProvider>
    </AuthProvider>
  );
}

export default App;
