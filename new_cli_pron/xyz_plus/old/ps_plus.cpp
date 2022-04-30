// #include <chrono>
// #include <thread>

#include "ps_plus.h"
//#include "data.h"

//using namespace std::chrono;

// #ifdef ASIO_STANDALONE
//   #include <asio.hpp>
//   #define ASIO asio
//   #define ASIO_ERROR_CODE asio::error_code
//   #define ASIO_TIMER asio::steady_timer
// #else
//   #include <boost/asio.hpp>
//   #define ASIO boost::asio
//   #define ASIO_ERROR_CODE boost::system::error_code
//   #define ASIO_TIMER boost::asio::deadline_timer
// #endif



// void sequential_encapsulated(PS_Data data[5]) {
//     XYZ_PocketSphinx ps[5];

// 	high_resolution_clock::time_point start;
// 	high_resolution_clock::time_point end;
//     high_resolution_clock::time_point starts[5];
// 	high_resolution_clock::time_point ends[5];

//     start=high_resolution_clock::now();
//     for(int i = 0; i< 5; i++) {
//         starts[i]=high_resolution_clock::now();
//         ps[i].init((void*)(data[i].jsgf.buffer), (size_t)data[i].jsgf.size, (void*)data[i].wav.buffer, (size_t)data[i].wav.size, (int)data[i].params.size, data[i].params.p);
//         ps[i].init_recognition();
//         ps[i].recognize_from_buffered_file();
//         ps[i].terminate();
//         ends[i]=high_resolution_clock::now();

//         strcpy(data[i].result.result,ps[i]._result);
//         data[i].result.size = ps[i]._result_size;
//     }
//     end=high_resolution_clock::now();

//     double total=0;
//     for(int i = 0; i< 5; i++) {
//         int frate = atoi(get_value(data[i].params.p,"-frate"));
//         printf("%d\t%s\t%lfms\n",frate,data[i].result.result,duration<double, std::milli>(ends[i] - starts[i]).count());
//         total = total + duration<double, std::milli>(ends[i] - starts[i]).count();
//     }

// 	auto dur_us = duration<double, std::micro>(end - start).count();
// 	auto dur_ms = duration<double, std::milli>(end - start).count();
// 	printf("Time: %lfus %lfms\n", dur_us, dur_ms);
//     printf("Computation time on task: %lfms\n", total);
// }

// void process(XYZ_PocketSphinx *ps, PS_Data *data) {
//         ps->init((void*)(data->jsgf.buffer), (size_t)data->jsgf.size, (void*)data->wav.buffer, (size_t)data->wav.size, (int)data->params.size, data->params.p);
//         ps->init_recognition();
//         ps->recognize_from_buffered_file();
//         ps->terminate(); 
// }

// void init(XYZ_PocketSphinx *ps, PS_Data *data) {
//         ps->init((void*)(data->jsgf.buffer), (size_t)data->jsgf.size, (void*)data->wav.buffer, (size_t)data->wav.size, (int)data->params.size, data->params.p);
//         ps->init_recognition();
// }

// void process_no_init(XYZ_PocketSphinx *ps) {
//        // ps->init((void*)(data->jsgf.buffer), (size_t)data->jsgf.size, (void*)data->wav.buffer, (size_t)data->wav.size, (int)data->params.size, data->params.p);
//        // ps->init_recognition();
//         ps->recognize_from_buffered_file();
//         ps->terminate(); 
// }

// void parallel_encapsulated(PS_Data data[5]) {
//     XYZ_PocketSphinx ps[5];
//     std::thread threads[5];

// 	high_resolution_clock::time_point start;
// 	high_resolution_clock::time_point end;
//     high_resolution_clock::time_point starts[5];
// 	high_resolution_clock::time_point ends[5];

//     start=high_resolution_clock::now();
//     // for(int i = 0; i< 5; i++) {
//     //     init(&ps[i], &data[i]);

//     // }
//     // printf("debug:");
//     // for(int i = 0; i< 5; i++) {
        
//     //     threads[i](process, &ps[i], &data[i]);
//     //     //process(&ps[i],&data[i]);
//     //     //process_no_init(&ps[i]);
        
//     // }
//     // for(int i = 0; i< 5; i++) {
//     //     threads[i].join();
//     // }

//     std::thread t0(process, &ps[0],&data[0]);
//     std::thread t1(process, &ps[1],&data[1]);
//     std::thread t2(process, &ps[2],&data[2]);
//     std::thread t3(process, &ps[3],&data[3]);
//     std::thread t4(process, &ps[4],&data[4]);
    
//     // int i=0;
//     // ps[i].init((void*)(data[i].jsgf.buffer), (size_t)data[i].jsgf.size, (void*)data[i].wav.buffer, (size_t)data[i].wav.size, (int)data[i].params.size, data[i].params.p);
//     // ps[i].init_recognition();
//     // ps[i].recognize_from_buffered_file();
//     // ps[i].terminate();
//     // std::thread t0(process_no_init, &ps[i]);
    
//     t0.join();
//     t1.join();
//     t2.join();
//     t3.join();
//     t4.join();
        


//     end=high_resolution_clock::now();

//     for(int i = 0; i< 5; i++) {
//         strcpy(data[i].result.result,ps[i]._result);
//         data[i].result.size = ps[i]._result_size;
//     }

//     double total=0;
//     for(int i = 0; i< 5; i++) {
//         int frate = atoi(get_value(data[i].params.p,"-frate"));
//         printf("%d\t%s\t\t%lfms\n",frate,data[i].result.result,duration<double, std::milli>(ends[i] - starts[i]).count());
//         total = total + duration<double, std::milli>(ends[i] - starts[i]).count();
//     }

// 	auto dur_us = duration<double, std::micro>(end - start).count();
// 	auto dur_ms = duration<double, std::milli>(end - start).count();
// 	printf("Time: %lfus %lfms\n", dur_us, dur_ms);
//     printf("Computation time on task: %lfms\n", total);
// }


int ps_plus_call(void* jsgf_buffer, int jsgf_buffer_size, void* audio_buffer, int audio_buffer_size, int argc, char *argv[], char* result, int rsize){


    //C version:
    //resultsize=ps_call_from_go(jsgf_buffer, (size_t)jsgf_buffer_size, audio_buffer, (size_t)audio_buffer_size, argc, argv, sresult);

    //Encapsulated version:
    XYZ_PocketSphinx ps;
    ps.init(jsgf_buffer, jsgf_buffer_size, audio_buffer, audio_buffer_size, argc, argv);
    ps.init_recognition();
    ps.recognize_from_buffered_file();
    ps.terminate();


    if (ps._result_size < rsize && strlen(ps._result)>0){

        for(int i=0;i<ps._result_size; i++){
            result[i]=(char)ps._result[i];
        }
    } 

    return ps._result_size;
 } 


// int
// main()
// {
//     struct PS_Data data[5];
//     load_data(params125, params125_size,
//               params72, params72_size, 
//               params80, params80_size,
//               params91, params91_size,
//               params105, params105_size, 
//               data);
//     //load_data(params125, params125, params125, params125, params125, data);


//     // DATA data(5);
//     // data.load(params125);

//     //sequential_encapsulated(data);
//     parallel_encapsulated(data);

//     // int i=4;
//     // XYZ_PocketSphinx ps1;
//     // ps1.init((void*)(data[i].jsgf.buffer), (size_t)data[i].jsgf.size, (void*)data[i].wav.buffer, (size_t)data[i].wav.size, (int)data[i].params.size, data[i].params.p);
//     // ps1.init_recognition();
//     // ps1.recognize_from_buffered_file();
//     // printf("%s\n",ps1._result);
//     // ps1.terminate();
    
//     printf("Working on it...\n");

//     return 0;

    
// }


