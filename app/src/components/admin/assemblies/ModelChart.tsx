import { Bar } from "react-chartjs-2";
import { useMemo } from "react";
import { useTranslation } from "react-i18next";
import { getModelColor } from "@/admin/colors.ts";

type ModelChartProps = {
  labels: string[];
  datasets: {
    model: string;
    data: number[];
  }[];
  dark?: boolean;
};
function ModelChart({ labels, datasets, dark }: ModelChartProps) {
  const { t } = useTranslation();
  const data = useMemo(() => {
    return {
      labels,
      datasets: datasets.map((dataset) => {
        return {
          label: dataset.model,
          data: dataset.data,
          backgroundColor: getModelColor(dataset.model),
        };
      }),
    };
  }, [labels, datasets]);

  const options = useMemo(() => {
    const text = dark ? "#fff" : "#000";

    return {
      scales: {
        x: {
          stacked: true,
          grid: {
            drawBorder: false,
            display: false,
          },
        },
        y: {
          beginAtZero: true,
          stacked: true,
          grid: {
            drawBorder: false,
            display: false,
          },
        },
      },
      plugins: {
        title: {
          display: false,
        },
        legend: {
          display: true,
          labels: {
            color: text,
          },
        },
      },
      color: text,
      borderWidth: 0,
      defaultFontColor: text,
      defaultFontSize: 16,
      defaultFontFamily: "Andika",
    };
  }, [dark]);

  return (
    <>
      <p className={`mb-2`}>{t("admin.model-chart")}</p>
      <Bar id={`model-chart`} data={data} options={options} />
    </>
  );
}

export default ModelChart;
