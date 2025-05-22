import pandas as pd
import matplotlib.pyplot as plt

def create_hist():
    df = pd.read_csv('perfomance-result/mean_throughput.csv')
    plt.bar(df['config'], df['mean_req_per_s'])
    plt.xlabel('Конфигурация')
    plt.ylabel('Средний throughput, req/s')
    plt.title('Сравнение средней пропускной способности')
    plt.grid(True)
    plt.show()
if __name__ == "__main__":
    create_hist()