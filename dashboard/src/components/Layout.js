import {
    BarChart2,
    Database,
    Home,
    Key,
    Settings,
    Terminal
} from 'lucide-react';
import Link from 'next/link';
import { useState } from 'react';
import { Menu, MenuItem, Sidebar } from 'react-pro-sidebar';

export default function Layout({ children }) {
  const [collapsed, setCollapsed] = useState(false);

  return (
    <div className="flex h-screen bg-gray-900">
      <Sidebar collapsed={collapsed} className="border-r border-gray-800">
        <Menu className="p-4">
          <MenuItem icon={<Home />} component={<Link href="/dashboard" />}>
            Dashboard
          </MenuItem>
          <MenuItem icon={<Terminal />} component={<Link href="/dashboard/playground" />}>
            AI Playground
          </MenuItem>
          <MenuItem icon={<Key />} component={<Link href="/dashboard/api-keys" />}>
            API Keys
          </MenuItem>
          <MenuItem icon={<Database />} component={<Link href="/dashboard/providers" />}>
            Providers
          </MenuItem>
          <MenuItem icon={<BarChart2 />} component={<Link href="/dashboard/analytics" />}>
            Analytics
          </MenuItem>
          <MenuItem icon={<Settings />} component={<Link href="/dashboard/settings" />}>
            Settings
          </MenuItem>
        </Menu>
      </Sidebar>

      <main className="flex-1 overflow-auto bg-gray-900 text-white">
        <nav className="bg-gray-800/50 backdrop-blur-sm border-b border-gray-800 p-4">
          <div className="flex items-center justify-between">
            <h1 className="text-xl font-bold">OrcaAI Dashboard</h1>
            <div className="flex items-center space-x-4">
              <button className="px-4 py-2 bg-blue-600 rounded-lg hover:bg-blue-700">
                New Query
              </button>
            </div>
          </div>
        </nav>
        
        <div className="p-6">
          {children}
        </div>
      </main>
    </div>
  );
}
