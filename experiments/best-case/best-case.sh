mode=$1

if [[ "${mode}" == "LAN" ]]
then
  echo "Running best-case LAN"
  # LAN experiment
  python3 experiments/best-case/test-automation.py LAN 3
  python3 experiments/best-case/summary.py LAN 3
fi

if [[ "${mode}" == "WAN" ]]
then
  echo "Running best-case WAN"
  # WAN experiment
  python3 experiments/best-case/test-automation.py WAN 3
  python3 experiments/best-case/summary.py WAN 3
  echo "Best case experiments done"
fi