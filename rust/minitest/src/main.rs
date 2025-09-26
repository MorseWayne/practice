use std::{
    process::id,
    sync::{Arc, Mutex},
    thread,
};

fn main() {
    let data = Arc::new(Mutex::new(Vec::new()));
    let mut handles = Vec::new();
    for i in 0..4 {
        let data1 = data.clone();
        handles.push(thread::spawn(move || {
            println!("{:?}", thread::current().id());
            data1.lock().unwrap().push(i);
        }));
    }
    for handle in handles {
        handle.join().unwrap();
    }

    let var = 1;
    let mut1 = &var;
    let mut2 = &var;

    println!("{:?}", data);
}
