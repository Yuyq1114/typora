import os
import random
import subprocess


def gpu_name() -> str:
    return subprocess.check_output(
        ["nvidia-smi", "--query-gpu=name", "--format=csv,noheader"],
        text=True,
    ).strip()


def main() -> None:
    stage = os.getenv("STAGE", "unknown")
    print(f"stage={stage}")
    print(f"gpu={gpu_name()}")

    # 训练内容故意保持简单，把注意力放在 Kubeflow 控制链路上。
    random.seed(7)
    xs = [(-1 + 2 * index / 999) for index in range(1_000)]
    ys = [4 * x - 1 + random.gauss(0, 0.04) for x in xs]

    weight = 0.0
    bias = 0.0
    learning_rate = 0.1

    for epoch in range(301):
        predictions = [weight * x + bias for x in xs]
        errors = [prediction - target for prediction, target in zip(predictions, ys)]
        loss = sum(error * error for error in errors) / len(errors)

        weight_gradient = 2 * sum(error * x for error, x in zip(errors, xs)) / len(xs)
        bias_gradient = 2 * sum(errors) / len(errors)
        weight -= learning_rate * weight_gradient
        bias -= learning_rate * bias_gradient

        if epoch % 100 == 0:
            print(f"epoch={epoch}, loss={loss:.6f}")

    print(f"model=y={weight:.4f}x{bias:+.4f}")
    print("training=complete")


if __name__ == "__main__":
    main()
