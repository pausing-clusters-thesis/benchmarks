import os
import json
import pandas as pd
import plotly.express as px
import plotly.io as pio

def load_data(root_dir):
    data = []

    for cluster in os.listdir(root_dir):
        cluster_path = os.path.join(root_dir, cluster)
        if not os.path.isdir(cluster_path):
            continue

        for driver in os.listdir(cluster_path):
            driver_path = os.path.join(root_dir, cluster, driver)
            if not os.path.isdir(driver_path):
                continue

            policy_path = os.path.join(driver_path, "pullPolicyNever")
            if not os.path.isdir(policy_path):
                continue

            for scenario in ["baseline", "busywait", "sidecar"]:
                scenario_path = os.path.join(policy_path, scenario)
                if not os.path.exists(scenario_path):
                    continue

                with open(scenario_path, "r") as f:
                    for line in f:
                        record = json.loads(line)
                        data.append({
                            "cluster": cluster,
                            "driver": driver,
                            "scenario": scenario,
                            "elapsed_time_ms": record["elapsed_time_ms"]
                        })

    return pd.DataFrame(data)

def plot_data(df, output_dir="plots"):
    os.makedirs(output_dir, exist_ok=True)

    grouped = df.groupby(["cluster", "driver"])

    cluster_dict = {
        "eks": "EKS",
        "gke": "GKE",
    }

    driver_dict = {
        "local.csi.scylladb.com": "local.csi.scylladb.com",
        "pd.csi.storage.gke.io": "pd.csi.storage.gke.io (pd-ssd)",
        "ebs.csi.aws.com": "ebs.csi.aws.com (gp3)"
    }

    for (cluster, driver), subset in grouped:
        fig = px.box(subset, x="scenario", y="elapsed_time_ms", color="scenario", color_discrete_map={"baseline": 'rgb(7,40,89)', "busywait": "rgb(9,56,125)", "sidecar": "rgb(8,81,156)"},
                     title=f"Time to reach Pod readiness in {cluster_dict[cluster]} with {driver_dict[driver]}, n=30",
                     labels={"elapsed_time_ms": "Elapsed Time (ms)", "scenario": "Scenario"})

        fig.update_traces(boxpoints='suspectedoutliers',
                          )

        fig.update_layout(showlegend=False,
                          )

        plot_filename = f"{output_dir}/{cluster}_{driver}_pullPolicyNever.png"
        pio.write_image(fig, plot_filename)
        print(f"Saved plot: {plot_filename}")

if __name__ == "__main__":
    root_directory = "./"
    df = load_data(root_directory)
    if not df.empty:
        plot_data(df)
    else:
        print("No data found!")
