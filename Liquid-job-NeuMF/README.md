# yao-NeuMF

```bash
pip install requests pandas typing tensorflow-gpu>=2.0

python official/recommendation/movielens.py --data_dir data/ --dataset ml-20m

# estimate time: 54m44s, Util ~50%
PYTHONPATH="$PYTHONPATH:/workspace" python /workspace/official/recommendation/ncf_keras_main.py \
--data_dir=/workspace/data/ \
--dataset=ml-20m \
--batch_size=1000
```