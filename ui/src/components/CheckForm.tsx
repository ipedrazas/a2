import { useState, FormEvent } from 'react';
import type { CheckRequest } from '../types';

interface CheckFormProps {
  onSubmit: (request: CheckRequest) => void;
  loading: boolean;
}

export function CheckForm({ onSubmit, loading }: CheckFormProps) {
  const [url, setUrl] = useState('');
  const [profile, setProfile] = useState('');
  const [target, setTarget] = useState('');
  const [skipChecks, setSkipChecks] = useState('');
  const [verbose, setVerbose] = useState(false);

  const handleSubmit = (e: FormEvent) => {
    e.preventDefault();

    const request: CheckRequest = {
      url: url.trim(),
      verbose: verbose,
    };

    if (profile) request.profile = profile;
    if (target) request.target = target;
    if (skipChecks.trim()) {
      request.skip_checks = skipChecks
        .split(',')
        .map(s => s.trim())
        .filter(s => s);
    }

    onSubmit(request);
  };

  const isValid = () => {
    return url.trim().length > 0 && !loading;
  };

  return (
    <form onSubmit={handleSubmit} className="space-y-4">
      <div>
        <label htmlFor="url" className="block text-sm font-medium text-gray-700 mb-1">
          GitHub URL
        </label>
        <input
          id="url"
          type="text"
          value={url}
          onChange={(e) => setUrl(e.target.value)}
          placeholder="https://github.com/owner/repo"
          className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
          disabled={loading}
        />
        <p className="mt-1 text-xs text-gray-500">
          Enter a GitHub repository URL
        </p>
      </div>

      <div className="grid grid-cols-2 gap-4">
        <div>
          <label htmlFor="profile" className="block text-sm font-medium text-gray-700 mb-1">
            Profile
          </label>
          <select
            id="profile"
            value={profile}
            onChange={(e) => setProfile(e.target.value)}
            className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
            disabled={loading}
          >
            <option value="">Auto-detect</option>
            <option value="cli">CLI</option>
            <option value="api">API</option>
            <option value="library">Library</option>
            <option value="desktop">Desktop</option>
          </select>
        </div>

        <div>
          <label htmlFor="target" className="block text-sm font-medium text-gray-700 mb-1">
            Target
          </label>
          <select
            id="target"
            value={target}
            onChange={(e) => setTarget(e.target.value)}
            className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
            disabled={loading}
          >
            <option value="">Auto-detect</option>
            <option value="poc">PoC</option>
            <option value="production">Production</option>
          </select>
        </div>
      </div>

      <div>
        <label htmlFor="skip" className="block text-sm font-medium text-gray-700 mb-1">
          Skip Checks (optional)
        </label>
        <input
          id="skip"
          type="text"
          value={skipChecks}
          onChange={(e) => setSkipChecks(e.target.value)}
          placeholder="go:race, common:k8s"
          className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
          disabled={loading}
        />
        <p className="mt-1 text-xs text-gray-500">
          Comma-separated check IDs to skip
        </p>
      </div>

      <div className="flex items-center">
        <input
          id="verbose"
          type="checkbox"
          checked={verbose}
          onChange={(e) => setVerbose(e.target.checked)}
          className="h-4 w-4 text-blue-600 focus:ring-blue-500 border-gray-300 rounded"
          disabled={loading}
        />
        <label htmlFor="verbose" className="ml-2 text-sm text-gray-700">
          Verbose output (show command output for failed/warning checks)
        </label>
      </div>

      <button
        type="submit"
        disabled={!isValid()}
        className="w-full px-6 py-3 bg-blue-600 text-white font-medium rounded-lg hover:bg-blue-700 disabled:bg-gray-400 disabled:cursor-not-allowed transition-colors"
      >
        {loading ? 'Submitting...' : 'Run Checks'}
      </button>
    </form>
  );
}
