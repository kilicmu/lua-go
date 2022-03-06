// let line;
// while (line = readline()) {
//     const [target, source] = line.split(" ");
//     let i = target.length - 1;
//     let j/*  */ = source.length - 1;
//     let idx = 0
//     while (i >= 0) {
//         idx = 0
//         while (j >= 0) {
//             if (target[i] === source[j]) {
//                 idx = j;
//                 j--;
//                 break
//             }
//             else {
//                 idx = j
//                 j--;
//                 continue
//             }
//         }
//         i--;
//     }
//     console.log(idx)
// }

// while (line = readline()) {
line = "3 7"
const sp = line.split(" ");
const c = parseInt(sp[0]), x = parseInt(sp[1]);
const nums = "3 4 7".split(" ").map((val) => parseInt(val)), length = c;

let result = 0, l = 0, r = 0, sum = nums[0];
while (l < length) {
    // 移动r直到超过x
    while (r < length && sum < x) {
        r++;
        sum += nums[r];
    }
    if (r === length && sum < x) {
        // r移动到结尾，还是没有超过x，则后续的循环都没有意义，break
        fin = true;
        break;
    }

    result += length - r;
    sum -= nums[l];
    l++;
}

console.log(result);
// }