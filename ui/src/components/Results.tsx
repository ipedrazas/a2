import type { JSONOutput } from '../types';

interface ResultsProps {
  result: JSONOutput | null;
}

const statusColors = {
  pass: 'bg-green-100 text-green-800 border-green-200',
  warn: 'bg-yellow-100 text-yellow-800 border-yellow-200',
  fail: 'bg-red-100 text-red-800 border-red-200',
  info: 'bg-gray-100 text-gray-800 border-gray-200',
};

export function Results({ result }: ResultsProps) {
  if (!result) {
    return null;
  }

  const { summary, maturity, results } = result;

  return (
    <div className="space-y-6">
      {/* Summary Card */}
      <div className="bg-white border border-gray-200 rounded-lg p-6 shadow-sm">
        <h2 className="text-xl font-semibold text-gray-800 mb-4">Summary</h2>

        <div className="grid grid-cols-2 md:grid-cols-4 gap-4 mb-6">
          <div className="text-center p-3 bg-blue-50 rounded-lg">
            <div className="text-2xl font-bold text-blue-600">{summary.total}</div>
            <div className="text-sm text-gray-600">Total Checks</div>
          </div>
          <div className="text-center p-3 bg-green-50 rounded-lg">
            <div className="text-2xl font-bold text-green-600">{summary.passed}</div>
            <div className="text-sm text-gray-600">Passed</div>
          </div>
          <div className="text-center p-3 bg-yellow-50 rounded-lg">
            <div className="text-2xl font-bold text-yellow-600">{summary.warnings}</div>
            <div className="text-sm text-gray-600">Warnings</div>
          </div>
          <div className="text-center p-3 bg-red-50 rounded-lg">
            <div className="text-2xl font-bold text-red-600">{summary.failed}</div>
            <div className="text-sm text-gray-600">Failed</div>
          </div>
        </div>

        <div className="flex items-center justify-between p-4 bg-gray-50 rounded-lg">
          <div>
            <div className="text-sm text-gray-600">Score</div>
            <div className="text-3xl font-bold text-gray-800">
              {summary.score.toFixed(1)}%
            </div>
          </div>
          <div>
            <div className="text-sm text-gray-600">Maturity Level</div>
            <div className="text-xl font-semibold text-gray-800">{maturity.level}</div>
          </div>
          <div className="text-right">
            <div className="text-sm text-gray-600">Duration</div>
            <div className="text-lg font-medium text-gray-800">
              {(summary.total_duration_ms / 1000).toFixed(2)}s
            </div>
          </div>
        </div>

        {maturity.description && (
          <div className="mt-4 p-3 bg-blue-50 border border-blue-200 rounded">
            <p className="text-sm text-blue-800">{maturity.description}</p>
          </div>
        )}

        {maturity.suggestions && maturity.suggestions.length > 0 && (
          <div className="mt-4">
            <h3 className="text-sm font-medium text-gray-700 mb-2">Suggestions</h3>
            <ul className="list-disc list-inside space-y-1">
              {maturity.suggestions.map((suggestion, idx) => (
                <li key={idx} className="text-sm text-gray-600">
                  {suggestion}
                </li>
              ))}
            </ul>
          </div>
        )}
      </div>

      {/* Results List */}
      <div className="bg-white border border-gray-200 rounded-lg p-6 shadow-sm">
        <h2 className="text-xl font-semibold text-gray-800 mb-4">Check Results</h2>

        <div className="space-y-3">
          {results.map((check) => (
            <div
              key={check.id}
              className={`p-4 border rounded-lg ${statusColors[check.status]}`}
            >
              <div className="flex items-start justify-between mb-2">
                <div>
                  <h3 className="font-medium text-gray-800">{check.name}</h3>
                  <p className="text-sm text-gray-600 font-mono">{check.id}</p>
                </div>
                <div className="text-right text-sm">
                  <div className="font-medium capitalize">{check.status}</div>
                  <div className="text-gray-600">{check.duration_ms}ms</div>
                </div>
              </div>

              {check.message && (
                <p className="text-sm mt-2 text-gray-700">{check.message}</p>
              )}

              {check.raw_output && (
                <details className="mt-2">
                  <summary className="text-sm font-medium cursor-pointer hover:text-gray-600">
                    Show output
                  </summary>
                  <pre className="mt-2 p-2 bg-gray-100 rounded text-xs overflow-x-auto">
                    {check.raw_output}
                  </pre>
                </details>
              )}
            </div>
          ))}
        </div>
      </div>
    </div>
  );
}
