import random
import subprocess


def main() -> None:
    gpu_name = subprocess.check_output(
        ["nvidia-smi", "--query-gpu=name", "--format=csv,noheader"],
        text=True,
    ).strip()
    print(f"GPU allocated by Kubernetes: {gpu_name}")

    random.seed(42)
    xs = [(-1 + 2 * index / 999) for index in range(1_000)]
    ys = [3 * x + 2 + random.gauss(0, 0.05) for x in xs]

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

        if epoch % 75 == 0:
            print(f"epoch={epoch:3d}, loss={loss:.6f}")

    print(f"result: y = {weight:.4f}x + {bias:.4f}")
    print("training completed")


if __name__ == "__main__":
    main()
