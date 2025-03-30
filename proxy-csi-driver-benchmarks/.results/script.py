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
                            "elapsed_time_s": record["elapsed_time_ms"]/1000
                        })

    return pd.DataFrame(data)

def plot_data(df, output_dir="plots"):
    os.makedirs(output_dir, exist_ok=True)

    driver_dict = {
        "local.csi.scylladb.com": "local.csi.scylladb.com",
        "pd.csi.storage.gke.io": "pd.csi.storage.gke.io (pd-ssd)",
        "ebs.csi.aws.com": "ebs.csi.aws.com (gp3)"
    }

    df['driver_label'] = df['driver'].map(driver_dict)
    grouped = df.groupby("cluster")

    cluster_dict = {
        "eks": "EKS",
        "gke": "GKE",
    }




    for cluster, subset in grouped:
        fig = px.box(subset, x="scenario", y="elapsed_time_s", color="driver_label",
                     color_discrete_map={"local.csi.scylladb.com": 'rgb(7,40,89)', "pd.csi.storage.gke.io (pd-ssd)": "rgb(9,56,125)", "ebs.csi.aws.com (gp3)": "rgb(8,81,156)"},
                     title=f"Time to reach Pod readiness in {cluster_dict[cluster]}, n=30",
                     category_orders={'scenario': ['baseline', 'busywait', 'sidecar'], 'driver_label': ['ebs.csi.aws.com (gp3)','pd.csi.storage.gke.io (pd-ssd)','local.csi.scylladb.com']},
                     labels={"elapsed_time_s": "Elapsed Time (s)", "scenario": "Scenario", "driver_label": "CSI storage provisioner:"},
                     boxmode='group',
                     )

        fig.update_traces(boxpoints='suspectedoutliers',
                          )

        fig.update_layout(showlegend=True,
                          legend=dict(
                              orientation="h",
                              yanchor="bottom",
                              y=1.02,
                              xanchor="right",
                              x=1
                          ),
                          )

        plot_filename = f"{output_dir}/{cluster}_pullPolicyNever.png"
        pio.write_image(fig, plot_filename)
        print(f"Saved plot: {plot_filename}")

if __name__ == "__main__":
    root_directory = "./"
    df = load_data(root_directory)
    if not df.empty:
        plot_data(df)
    else:
        print("No data found!")
