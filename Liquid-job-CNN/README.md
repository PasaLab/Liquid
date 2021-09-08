# yao-mnist

```bash
pip install requests tensorflow-gpu==1.14

# estimate time: 11min 30s
PYTHONPATH="$PYTHONPATH:/workspace" python /workspace/official/r1/mnist/mnist.py --data_dir=/workspace/data/

# estimate time: 3min 42s
PYTHONPATH="$PYTHONPATH:/workspace" python /workspace/official/r1/mnist/mnist_test.py --benchmarks=.
```