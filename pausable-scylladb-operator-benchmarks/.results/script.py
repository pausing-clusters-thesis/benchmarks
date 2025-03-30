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

            for nodes in os.listdir(driver_path):
                nodes_path = os.path.join(driver_path, nodes)
                if not os.path.isdir(nodes_path):
                    continue

                for scenario in ["baseline", "cold", "prewarmed"]:
                    scenario_path = os.path.join(nodes_path, scenario)
                    if not os.path.exists(scenario_path):
                        continue

                    with open(scenario_path, "r") as f:
                        for line in f:
                            record = json.loads(line)
                            data.append({
                                "cluster": cluster,
                                "driver": driver,
                                "nodes": nodes,
                                "scenario": scenario,
                                "component": "Application",
                                "elapsed_time_s": record["application_time_ms"]/1000
                            })

                            data.append({
                                "cluster": cluster,
                                "driver": driver,
                                "nodes": nodes,
                                "scenario": scenario,
                                "component": "Overhead",
                                "elapsed_time_s": record["overhead_time_ms"]/1000
                            })

    return pd.DataFrame(data)

def plot_data(df, output_dir="plots"):
    os.makedirs(output_dir, exist_ok=True)

    grouped = df.groupby(["cluster", "driver", "nodes"])

    cluster_dict = {
        "eks": "EKS",
        "gke": "GKE",
    }

    driver_dict = {
        "gke": "pd.csi.storage.gke.io (pd-ssd)",
        "eks": "ebs.csi.aws.com (gp3)"
    }

    for (cluster, driver, nodes), subset in grouped:
        fig = px.box(subset, x="scenario", y="elapsed_time_s", color="component",
                     color_discrete_map={"Overhead": "rgb(8,81,156)", "Application": "rgb(219, 64, 82)"},
                     title=f"Time to unpause ScyllaDB cluster in {cluster_dict[cluster]} with {driver_dict[cluster]}, n=30",
                     category_orders={'scenario': ['baseline', 'cold', 'prewarmed'], 'component': ['Overhead', 'Application']},
                     labels={"elapsed_time_s": "Time (s)", "scenario": "Scenario", "component": "Component"},
                     boxmode='group')

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

        plot_filename = f"{output_dir}/{cluster}_{driver}_{nodes}.png"
        pio.write_image(fig, plot_filename)
        print(f"Saved plot: {plot_filename}")

if __name__ == "__main__":
    root_directory = "./"
    df = load_data(root_directory)
    if not df.empty:
        plot_data(df)
    else:
        print("No data found!")
