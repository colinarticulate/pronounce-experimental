'''
Report. From a list of .toml files with results in it of test_pronounce, it yields an excel file
reporting the outcome of all results in one single spreadsheet
'''

import os
import shutil
import subprocess
import time
import datetime
import toml
import stat
import numpy as np
from collections import defaultdict
import xlsxwriter as xl

from sklearn.metrics import accuracy_score, f1_score, precision_score, recall_score, roc_auc_score, matthews_corrcoef

metrics={
    'Accuracy': accuracy_score,
    #'Precision': precision_score,
    #'Recall': recall_score,
    #'F1': f1_score,
    #'ROC AUC': roc_auc_score,
    #'MCC': matthews_corrcoef
}


def read_result(result_file, results_folder):
    
    data = toml.load(os.path.join(results_folder, result_file))

    info={}
    for item in data['info']:
        info[item]=data['info'][item]

    predictions={}
    for item in data['predictions']:
        predictions[item]=data['predictions'][item]

    results={}
    for item in data['results']:
        results[item]=data['results'][item]

    return info, predictions, results


def add_items(worksheet,predictions):
    items = list(predictions.keys())    
    worksheet.write_column(0,0, items)

def add_items_hyperlink_format(worksheet, predictions, path, file_extension):
    items = list(predictions.keys())
    for i, item in enumerate(items): 
        audio_file = "_".join(item.split("_")[:-1])   
        worksheet.write_url(i,0, os.path.join(path,f"{audio_file}.{file_extension}") ,string=item)


def add_predictions(worksheet, col, predictions):

    for i, prediction in enumerate(list(predictions.keys())):
        worksheet.write(i, col, predictions[prediction])

def format_result(result):
    #print(result)
    formatted_result = [" ".join(item) for item in result]
    return formatted_result


def add_results(worksheet, predictions, results, audios_folder):
    #add_items(worksheet, predictions)
    add_items_hyperlink_format(worksheet, predictions, audios_folder, "wav")
    add_predictions(worksheet, 1, predictions)
    for i, result in enumerate(results):
        formatted_results = format_result(results[result])
        worksheet.write_row(i, 2, formatted_results)


def conditional_formatting(workbook, worksheet, start_row, start_col, n_rows, n_cols):
    
    # Add a format. Light red fill with dark red text.
    format_fail = workbook.add_format({'bg_color': '#FFC7CE',
                                       'font_color': '#9C0006'})

    # Add a format. Light red fill with dark red text.
    format_blank = workbook.add_format({'bg_color': 'red',
                                        'font_color': 'red'})

    # Add a format. Green fill with dark green text.
    format_pass = workbook.add_format({'bg_color': '#C6EFCE',
                                       'font_color': '#006100'})


    worksheet.conditional_format(start_row, start_col, n_rows, n_cols, {'type':     'text',
                                                        'criteria': 'containing',
                                                        'value':    'PASS',
                                                        'format':   format_pass})

    worksheet.conditional_format(start_row, start_col, n_rows, n_cols, {'type':     'text',
                                                        'criteria': 'containing',
                                                        'value':    'FAIL',
                                                        'format':   format_fail})

    worksheet.conditional_format(start_row, start_col, n_rows, n_cols, {'type':     'blanks',
                                                       # 'criteria': 'containing',
                                                       # 'value':    '',
                                                        'format':   format_blank})


def conditional_formatting_on_results(workbook, worksheet, start_row, start_col, n_rows, n_cols):
    
    # Add a format. Light red fill with dark red text.
    format_missing = workbook.add_format({'bg_color': '#6BEBEF',
                                       'font_color': '#063637'})

    worksheet.conditional_format(start_row, start_col, n_rows, n_cols, {'type':     'text',
                                                        'criteria': 'containing',
                                                        'value':    'missing',
                                                        'format':   format_missing})

        # Add a format. Light red fill with dark red text.
    format_surprise = workbook.add_format({'bg_color': '#FBB637',
                                       'font_color': '#503502'})

    worksheet.conditional_format(start_row, start_col, n_rows, n_cols, {'type':     'text',
                                                        'criteria': 'containing',
                                                        'value':    'surprise',
                                                        'format':   format_surprise})



def transform_to_numpy(predictions):
    np_predictions = np.array(list(predictions.values()))
    #fail_set = np.argwhere(np_predictions=='FAIL')[:,0]
    pass_set = np.argwhere(np_predictions=='PASS')[:,0]
    #empty_set = np.argwhere(np_predictions=='')[:,0]

    y_pred = np.zeros(np_predictions.shape).astype(int)
    y_pred[pass_set] = 1

    y_true = np.ones(np_predictions.shape).astype(int)

    return y_true, y_pred


def calculate_metrics(predictions):
    y_true, y_pred = transform_to_numpy(predictions)

    scores={}
    for metric in metrics:
        score = 100*metrics[metric](y_true, y_pred)
        scores[metric]= f"{score:.2f}%"

    return scores


def add_metrics(worksheet, row, col, predictions):
    
    scores = calculate_metrics(predictions)

    for i,score in enumerate(scores):
        worksheet.write(row+i, col, scores[score])

    return list(scores.keys())


def add_model_info(workbook, worksheet, col, worksheet_legend, info, predictions, results):
    row = len(predictions)

    info_column=[]
    #Legend
    worksheet_legend.write_string(col,0,f"{col}")
    cell_format=workbook.add_format()
    cell_format.set_bold(True)
    worksheet_legend.write_string(col,1, info['model_name'])
    worksheet_legend.set_column(col, 1, len(info['model_name']), cell_format)
    info_column.append("Model (see Legend sheet)")

    #Results
    worksheet.write_string(row, col, f"{col}") #info['model_name'])

    metrics_keys = add_metrics(worksheet, row+1, col, predictions)
    info_column=info_column+metrics_keys
    #here we could add more info

    return info_column


def add_info(workbook, worksheet, worksheet_legend, models_info):

    model_names=[]
    for i, model_info in enumerate(models_info):

        info, predictions = model_info
        model_names.append(info['model_name'])
        row=len(predictions)
        col=i+1

        info_column=[]
        #Legend
        worksheet_legend.write_string(col,0,f"{col}")
        cell_format=workbook.add_format()
        cell_format.set_bold(True)
        worksheet_legend.write_string(col,1, info['model_name'], cell_format)
        #worksheet_legend.set_column(col, 1, len(info['model_name']), cell_format)
        info_column.append("Model (see Legend sheet)")

        #Results
        worksheet.write_string(row, col, f"{col}") #info['model_name'])

        metrics_keys = add_metrics(worksheet, row+1, col, predictions)
        info_column=info_column+metrics_keys
        #here we could add more info

    right_align_format = workbook.add_format()
    right_align_format.set_align('right')
    right_align_format.set_bold(True)
    worksheet.write_column(row, 0, info_column, right_align_format)

    max_len = calculate_max_len(model_names)
    worksheet_legend.set_column(1,1,max_len+5)


def calculate_max_len(items):
    
    max_len=0
    for prediction in items:
        n = len(prediction)
        if n> max_len:
            max_len=n

    return max_len


def formatting_report(workbook, worksheet, worksheet_legend, nrow, ncol, predictions):

    max_len = calculate_max_len(predictions)

    worksheet.set_column(0,0,max_len+5)

    contents_format = workbook.add_format()
    contents_format.set_align('center')
    #

    worksheet.set_column(1,200, None, contents_format)
   

def format_results_sheet(workbook, worksheet, predictions):

    max_len = calculate_max_len(predictions)

    worksheet.set_column(0,0,max_len+5)
    worksheet.set_column(2,200,12) #the lenght of someting like: "ehr missing"
    

def create_report(report_file, experiment_name, results_files, results_folder):

    workbook = xl.Workbook(report_file)
    worksheet = workbook.add_worksheet("Results") #This could be timestamps as well
    worksheet_legend = workbook.add_worksheet("Legend")

    model_info=[]
    for i,file in enumerate(results_files):
        info, predictions, results = read_result(file, results_folder)
        if i == 0:
            add_items(worksheet, predictions)

        add_predictions(worksheet, i+1, predictions)
        worksheet_i = workbook.add_worksheet(f"{i+1}")
        add_results(worksheet_i, predictions, results, info['audios_folder'])
        format_results_sheet(workbook, worksheet_i, predictions)
        conditional_formatting_on_results(workbook, worksheet_i, 0, 2, len(results), 300)
        conditional_formatting(workbook, worksheet_i, 0, 1, len(predictions)-1, 1)
        model_info.append((info,predictions))
        
        
    #add_model_info(workbook, worksheet, i+1, worksheet_legend, info, predictions, results)

    add_info(workbook, worksheet, worksheet_legend, model_info)
    formatting_report(workbook, worksheet, worksheet_legend, len(predictions)-1, i+1, predictions )
    conditional_formatting(workbook, worksheet, 0, 1, len(predictions)-1, i+1)


    workbook.close()

def main():


    report_dir = "./../Reports"
    experiment_name = "Data_augmentation_testing"
    report_file = os.path.join(report_dir, f"{experiment_name}.xlsx")
    results_folder = "./../Results"

    results_files = [f for f in os.listdir(results_folder) if f.endswith(".toml")]
    #results_files = ['Bare_loudness_speed_x_Test_Harness.toml']
    

    create_report(report_file, experiment_name, results_files, results_folder)
    
    print("finished.main")

if __name__=='__main__':
    start=time.time()
    main()
    stop=time.time()
    print("Finished.")
    print(f"Time: {stop-start} seconds.")