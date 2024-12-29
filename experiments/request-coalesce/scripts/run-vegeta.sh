vegeta attack -targets vegeta -timeout=10s -duration=30s | tee stress-test-results/results.bin | vegeta report
  vegeta report -type=json stress-test-results/results.bin > stress-test-results/metrics.json
  cat stress-test-results/results.bin | vegeta plot > stress-test-results/plot.html
  cat stress-test-results/results.bin | vegeta report -type="hist[0,100ms,200ms,300ms,500ms,1000ms,2000ms]"